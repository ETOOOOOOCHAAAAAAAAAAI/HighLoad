package main

import (
	"log"
	"net/http"

	"Image-Hosting/internal/handler"
	"Image-Hosting/internal/repository"
	"Image-Hosting/internal/service"
)

func main() {
	pg := repository.NewPostgresRepo()
	redis := repository.NewRedisCache()

	imgService := service.NewImageService(pg, redis)
	imgHandler := handler.NewImageHandler(imgService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /images", imgHandler.CreateImage)
	mux.HandleFunc("PUT /s3/upload/{id}", imgHandler.UploadToS3)
	mux.HandleFunc("GET /images/{id}", imgHandler.GetImage)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
