package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"
)

func main() {
	router := gin.Default()

	var store service.Storage

	storageType := os.Getenv("STORAGE")

	if storageType == "postgres" {
		conn := os.Getenv("POSTGRES_URL")
		if conn == "" {
			log.Fatal("POSTGRES_URL is not set")
		}
		pg, err := postgres.New(conn)
		if err != nil {
			log.Fatal(err)
		}

		store = pg
	} else {
		store = memory.New()
	}

	svc := service.New(store)
	h := handler.New(svc)

	router.POST("/shorten", h.Create)
	router.GET("/:code", h.Get)

	log.Println("server started on :8080")

	router.Run(":8080")
}
