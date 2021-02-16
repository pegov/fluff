package db

import (
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
	GetKey() string
	SetLink(link.Link) bool
}

// Database ...
type Database struct {
	*redis.Client
	logger            *zap.SugaredLogger
	queue             chan string
	resource          int
	linkDefaultLen    int
	defaultExpiration time.Duration
}

var ctx = context.Background()

// NewDatabase ...
func NewDatabase(url string, logger *zap.SugaredLogger) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})
	// DEV
	rdb.FlushAll(ctx)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	queue := make(chan string, 0)
	db := &Database{
		rdb,
		logger,
		queue,
		500,
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

// SetLink ...
func (db *Database) SetLink(link link.Link) bool {
	ok, _ := db.SetNX(ctx, link.Key, link.URL, db.defaultExpiration).Result()
	return ok
}

func (db *Database) isFree(key string) bool {
	err := db.Get(ctx, key).Err()
	return err == redis.Nil
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func (db *Database) createKey() {
	b := make([]byte, db.linkDefaultLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	db.queue <- string(b)
}

// GetKey ...
func (db *Database) GetKey() string {
	var key string
	for {
		key = <-db.queue
		go db.createKey()
		if db.isFree(key) {
			return key
		}
	}
}

func (db *Database) createInitialKeys() {
	for i := 0; i < db.resource; i++ {
		go db.createKey()
	}
}
