package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	. "github.com/pegov/fluff/model"
	"github.com/pegov/fluff/util"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type LinkRepo struct {
	*sqlx.DB
	cache *redis.Client
}

func NewLinkRepo(db *sqlx.DB, cache *redis.Client) *LinkRepo {
	return &LinkRepo{
		db,
		cache,
	}
}

var ctx = context.Background()

func (db *LinkRepo) Create(long string, shortLen int) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	var short string
	tries := 0
	for {
		tries = tries + 1
		short = util.RandomString(shortLen)
		err = tx.QueryRow("SELECT id FROM link WHERE short = $1", short).Err()
		if err == sql.ErrNoRows {
			break
		} else if err == nil {
			if tries < 500 {

				tries++
				continue
			} else {
				return "", errors.New("tries > 500")
			}
		} else {
			return "", err
		}
	}

	_, err = tx.Exec("INSERT INTO link(short, long) VALUES ($1, $2)", short, long)
	if err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}

	return short, nil
}

func (db *LinkRepo) GetAllLinks() ([]Link, error) {
	links := []Link{}
	err := db.Select(&links, "SELECT id, short, long FROM link")
	if err == nil || err == sql.ErrNoRows {
		return links, nil
	} else {
		return links, err
	}
}

func (db *LinkRepo) GetByShort(short string) (*Link, error) {
	var link Link
	key := fmt.Sprintf("link:short:%v", short)

	val, err := db.cache.Get(ctx, key).Result()
	switch err {
	case redis.Nil:
		err := db.QueryRowx("SELECT id, short, long FROM link WHERE short = $1", short).StructScan(&link)
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		case nil:
			buf, err := json.Marshal(link)
			if err != nil {
				return nil, err
			}
			err = db.cache.SetEx(ctx, key, buf, 60*time.Second).Err()
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return &link, nil
		default:
			return nil, err
		}
	case nil:
		err = json.Unmarshal([]byte(val), &link)
		if err != nil {
			return nil, err
		}
		return &link, nil
	default:
		return nil, err
	}
}

func (db *LinkRepo) DeleteByShort(short string) error {
	_, err := db.Exec("DELETE FROM link WHERE short = $1", short)
	if err != nil {
		return err
	}

	err = db.cache.Del(ctx, fmt.Sprintf("link:short:%v", short)).Err()
	if err != nil {
		return err
	}

	return nil
}
