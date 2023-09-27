package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pegov/fluff/api"
	"github.com/pegov/fluff/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "0.0.0.0"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	bindAddr := fmt.Sprintf("%v:%v", host, port)

	pgxDb, err := sqlx.Connect("pgx", "postgres://postgres:postgres@127.0.0.1:5432/fluff")

	if err != nil {
		log.Fatalln(err)
	}

	pgxDb.DB.SetMaxOpenConns(10)

	rOptions, err := redis.ParseURL("redis://127.0.0.1:6379/1")
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(rOptions)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln(err)
	}

	repo := db.NewLinkRepo(pgxDb, rdb)

	rand.Seed(time.Now().UnixNano())

	r := gin.Default()

	r.GET("/api/links", api.GetAllLinks(repo))
	r.POST("/api/links", api.CreateLink(repo))

	r.GET("/api/links/:short", api.GetLink(repo))
	r.DELETE("/api/links/:short", api.DeleteLink(repo))

	r.Run(bindAddr)
}
