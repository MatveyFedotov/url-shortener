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

	var store interface {
		Save(string, string) error
		Get(string) (string, error)
		FindByURL(string) (string, error)
		Exists(string) (bool, error)
	}

	storageType := os.Getenv("STORAGE")

	if storageType == "postgres" {
		conn := "postgres://postgres:postgres@postgres:5432/shortener"

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
