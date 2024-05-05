package service

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func ValidFullURL(fullURL string) error {
	if fullURL == "" {
		return fmt.Errorf("empty url")
	}

	pattern := `^(https?://|www.)?[a-zA-Z0-9-]{1,256}([.][a-zA-Z-]{1,256})?([.][a-zA-Z]{1,30})([/][a-zA-Z0-9/?=%&#_.-]+)`
	if valid, _ := regexp.Match(pattern, []byte(fullURL)); !valid {
		return fmt.Errorf("%s is a invalid url", fullURL)
	}
	return nil
}

func ValidShortURL(shortURL string) error {
	if shortURL == "" {
		return fmt.Errorf("empty url")
	}

	pattern := `^[a-zA-Z0-9_]{10}`
	if valid, _ := regexp.Match(pattern, []byte(shortURL)); !valid {
		return fmt.Errorf("%s is a invalid short url", shortURL)
	}
	return nil
}

func GenerateShortURL() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_")

	result := make([]rune, 10)
	for i := 0; i < 10; i++ {
		result[i] = alphabet[rnd.Intn(len(alphabet))]
	}
	return string(result)
}
