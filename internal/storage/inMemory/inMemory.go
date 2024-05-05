package inMemory

import (
	"fmt"
	"ozon-fintech/internal/models"
	"ozon-fintech/internal/storage"
	"sync"
)

type InMemory struct {
	mu             sync.Mutex
	shortToFullUrl map[string]string
	fullUrlToShort map[string]string
}

func New() *InMemory {
	return &InMemory{
		mu:             sync.Mutex{},
		shortToFullUrl: make(map[string]string),
		fullUrlToShort: make(map[string]string),
	}
}

func (i *InMemory) GetFullURL(shortUrl string) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	res, ok := i.shortToFullUrl[shortUrl]
	if !ok {
		return "", fmt.Errorf("there are no urls with this shortUrl %s", shortUrl)
	}

	return res, nil
}

func (i *InMemory) LoadShortURL(link models.Link) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.fullUrlToShort[link.FullUrl]; ok {
		return "", storage.DuplicateErr
	}

	i.fullUrlToShort[link.FullUrl] = link.ShortUrl
	i.shortToFullUrl[link.ShortUrl] = link.FullUrl
	return link.ShortUrl, nil
}

func (i *InMemory) Stop() error {
	return nil
}
