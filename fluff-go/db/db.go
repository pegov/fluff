package db

import (
	"container/list"
	"context"
	"math/rand"
	"time"

	"github.com/NoSoundLeR/fluff/fluff-go/link"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Getter ...
type Getter interface {
	GetLink(string) (string, error)
}

//Setter ...
type Setter interface {
	GetFreeKey() string
	SetLink(link.Link, bool) bool
}

// Database ...
type Database struct {
	*redis.Client
	logger            *zap.SugaredLogger
	queue             *list.List
	resource          int
	linkDefaultLen    int
	defaultExpiration time.Duration
}

var ctx = context.Background()

// NewDatabase ...
func NewDatabase(logger *zap.SugaredLogger) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	db := &Database{
		rdb,
		logger,
		list.New(),
		10000,
		6,
		time.Hour * 12,
	}
	db.createInitialKeys()
	return db
}

// GetLink ...
func (db *Database) GetLink(key string) (string, error) {
	url, err := db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	db.Del(ctx, key)
	return url, nil
}

func (db *Database) queueContainsValue(value string) bool {
	if db.queue.Len() == 0 {
		return false
	}
	elem := db.queue.Back()
	if elem.Value == value {
		return true
	}
	for {
		elem = elem.Next()
		if elem == nil {
			return false
		}
		if elem.Value == value {
			return true
		}
	}
}

func (db *Database) queueRemoveByValue(value string) bool {
	if db.queue.Len() == 0 {
		return false
	}
	elem := db.queue.Back()
	if elem.Value == value {
		db.queue.Remove(elem)
		return true
	}
	for {
		elem = elem.Next()
		if elem == nil {
			return false
		}
		if elem.Value == value {
			db.queue.Remove(elem)
			return true
		}
	}
}

// SetLink ...
func (db *Database) SetLink(link link.Link, custom bool) bool {
	// TODO
	ok, _ := db.SetNX(ctx, link.Key, link.URL, db.defaultExpiration).Result()
	if ok {
		var removed bool
		if custom {
			removed = db.queueRemoveByValue(link.Key)
		}
		if !custom || removed {
			go db.createFreeKey()
		}
	}
	return ok
}

func (db *Database) isFree(key string) bool {
	_, err := db.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return true
		}
	}
	return false
}

func (db *Database) createInitialKeys() {
	for i := 0; i < db.resource; i++ {
		db.createFreeKey()
	}
}

func (db *Database) createFreeKey() string {
	b := make([]byte, db.linkDefaultLen)
	for {
		key := createRandomKey(&b)
		if !db.queueContainsValue(key) && db.isFree(key) {
			db.queue.PushBack(key)
			return key
		}
	}
}

// GetFreeKey ...
func (db *Database) GetFreeKey() string {
	front := db.queue.Front()
	defer db.queue.Remove(front)
	return front.Value.(string)
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func createRandomKey(b *[]byte) string {
	for i := range *b {
		(*b)[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(*b)

}
