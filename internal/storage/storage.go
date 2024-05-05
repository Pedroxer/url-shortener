package storage

import (
	"errors"
	"log/slog"
	"ozon-fintech/internal/models"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go
type DbType interface {
	LoadShortURL(link models.Link) (string, error)
	GetFullURL(shortUrl string) (string, error)
	Stop() error
}

var DuplicateErr = errors.New("duplicated short url")

type Storage struct {
	Storage DbType
	logger  *slog.Logger
}

func New(storage DbType, logger *slog.Logger) *Storage {
	return &Storage{
		Storage: storage,
		logger:  logger,
	}
}

func (s *Storage) Stop() error {
	return s.Storage.Stop()
}
