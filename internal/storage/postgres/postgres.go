package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"ozon-fintech/internal/config"
	"ozon-fintech/internal/models"
	"ozon-fintech/internal/storage"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg *config.Postgres) (*Postgres, error) {
	DSN := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
		cfg.Database,
		cfg.User,
		cfg.Password,
		cfg.Addr,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ConnectToDb: %w", err)
	}

	return &Postgres{db: db}, nil

}

func (p *Postgres) GetFullURL(shortUrl string) (string, error) {
	query := `SELECT url FROM link WHERE short_url = $1`

	var res string

	if err := p.db.QueryRow(query, shortUrl).Scan(&res); err != nil {
		return "", err
	}

	return res, nil
}

func (p *Postgres) LoadShortURL(link models.Link) (string, error) {
	checkQuery := `SELECT count(*) FROM link WHERE short_url = $1`
	var amount int
	if err := p.db.QueryRow(checkQuery, link.ShortUrl).Scan(&amount); err != nil {
		if amount > 1 {
			return "", storage.DuplicateErr
		}
	}
	query := `INSERT INTO link (url, short_url) VALUES ($1, $2)`
	if _, err := p.db.Exec(query, link.FullUrl, link.ShortUrl); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"link_pkey\"" {
			return "", storage.DuplicateErr
		}
		return "", err

	}
	return link.ShortUrl, nil
}

func (p *Postgres) Stop() error {
	return p.db.Close()
}
