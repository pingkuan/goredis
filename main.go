package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Post struct {
	UserID int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Posts []Post

var ctx = context.Background()

func main() {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r.GET("/", func(c *gin.Context) {
		val, err := rdb.Get(ctx, "posts").Result()
		if err == redis.Nil {
			log.Println("posts does not exist")
			res, err := http.Get("https://jsonplaceholder.typicode.com/posts")
			if err != nil {
				log.Fatal(err)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			res.Body.Close()

			var posts Posts
			json.Unmarshal(body, &posts)

			rdb.SetEx(ctx, "posts", body, 5*time.Second)

			c.JSON(http.StatusOK, posts)
			return
		} else if err != nil {
			panic(err)
		}

		var posts Posts
		data := []byte(val)
		json.Unmarshal(data, &posts)
		c.JSON(http.StatusOK, posts)
	})

	r.Run(":8080")
}
