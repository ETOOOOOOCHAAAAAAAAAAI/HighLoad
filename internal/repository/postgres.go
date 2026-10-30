package repository

import (
	"errors"
	"sync"

	"Image-Hosting/internal/models"
)

var ErrNotFound = errors.New("image not found")

type PostgresRepo struct {
	mu   sync.RWMutex
	data map[string]models.Image
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{
		data: make(map[string]models.Image),
	}
}

func (r *PostgresRepo) Save(img models.Image) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[img.ID] = img
	return nil
}

func (r *PostgresRepo) UpdateStatus(id, status, cdnURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	img, exists := r.data[id]
	if !exists {
		return ErrNotFound
	}

	img.Status = status
	img.CDNUrl = cdnURL
	r.data[id] = img
	return nil
}

func (r *PostgresRepo) GetByID(id string) (models.Image, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	img, exists := r.data[id]
	if !exists {
		return models.Image{}, ErrNotFound
	}
	return img, nil
}
