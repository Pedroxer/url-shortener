package service

import (
	"ozon-fintech/internal/models"
	"ozon-fintech/internal/storage"
)

type Service struct {
	storage storage.DbType
}

func NewService(storage storage.DbType) *Service {
	return &Service{storage: storage}
}

//go:generate mockgen -source=shortener.go -destination=mocks/mock.go

type Services interface {
	GetFullURL(shortURL string) (string, error)
	LoadShortURL(link models.Link) (string, error)
}

func (s *Service) GetFullURL(shortURL string) (string, error) {
	fullURL, err := s.storage.GetFullURL(shortURL)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

func (s *Service) LoadShortURL(link models.Link) (string, error) {
	var err error
	link.ShortUrl = GenerateShortURL()
	link.ShortUrl, err = s.storage.LoadShortURL(link)
	if err != nil {
		return "", err
	}
	return link.ShortUrl, nil
}
