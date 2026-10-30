package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"Image-Hosting/internal/models"
	"Image-Hosting/internal/repository"
)

type ImageService struct {
	db    *repository.PostgresRepo
	cache *repository.RedisCache
	queue chan string
}

func NewImageService(db *repository.PostgresRepo, cache *repository.RedisCache) *ImageService {
	s := &ImageService{
		db:    db,
		cache: cache,
		queue: make(chan string, 100),
	}
	go s.startWorker()
	return s
}

func (s *ImageService) CreatePendingImage(filename string) (string, string, error) {
	id := fmt.Sprintf("img_%d", rand.Intn(100000))

	img := models.Image{
		ID:        id,
		Filename:  filename,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.db.Save(img); err != nil {
		return "", "", err
	}

	presignedURL := fmt.Sprintf("http://localhost:8080/s3/upload/%s", id)
	return id, presignedURL, nil
}

func (s *ImageService) ConfirmS3Upload(id string) {
	log.Printf("[S3] Файл %s загружен, отправляем событие в Queue", id)
	s.queue <- id
}

func (s *ImageService) startWorker() {
	for id := range s.queue {
		cdnURL := fmt.Sprintf("https://cdn.example.com/images/%s", id)
		if err := s.db.UpdateStatus(id, models.StatusReady, cdnURL); err != nil {
			log.Printf("[Worker] Ошибка обновления статуса для %s: %v", id, err)
			continue
		}
		s.cache.Delete(id)
		log.Printf("[Worker] Статус %s обновлен на READY", id)
	}
}

func (s *ImageService) GetMetadata(id string) (models.Image, string, error) {
	if img, ok := s.cache.Get(id); ok {
		return img, "HIT", nil
	}

	img, err := s.db.GetByID(id)
	if err != nil {
		return models.Image{}, "MISS", err
	}

	s.cache.Set(id, img)
	return img, "MISS", nil
}
