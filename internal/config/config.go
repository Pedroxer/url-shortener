package config

type Config struct {
	PgStorage Postgres
}

type Postgres struct {
	Addr     string
	Port     int
	Database string
}
