package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pegov/fluff/db"
	"github.com/pegov/fluff/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func SetupRouter() *gin.Engine {
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

	r := gin.Default()

	r.GET("/api/links", handler.GetAllLinks(repo))
	r.POST("/api/links", handler.CreateLink(repo))

	r.GET("/api/links/:short", handler.GetLink(repo))
	r.DELETE("/api/links/:short", handler.DeleteLink(repo))

	r.GET("/:short", handler.RedirectToLink(repo))
	return r
}

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

	rand.Seed(time.Now().UnixNano())

	r := SetupRouter()
	r.Run(bindAddr)
}
