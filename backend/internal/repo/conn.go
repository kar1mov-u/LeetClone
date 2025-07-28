// Creates a connection pool(pgx.Pool) to Database(postgres). For dev enviroment, runs migrations everytime server is started
package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST"`
	Port     uint16 `envconfig:"POSTGRES_PORT"`
	Username string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	Database string `envconfig:"POSTGRES_DB"`
	MaxConns int32  `envconfig:"POSTGRES_MAX_CONNS"`
	MinConns int32  `envconfig:"POSTGRES_MIN_CONNS"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

func (cfg Config) GetConnString() string {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	return connString
}

// create a connection pool to DB
func NewConn(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := cfg.GetConnString()
	connString += fmt.Sprintf("&pool_max_conns=%d&pool_min_conns=%d",
		cfg.MaxConns,
		cfg.MinConns,
	)
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return conn, err
}

func Migrate(ctx context.Context, cfg Config, migrationsPath string) error {
	connString := cfg.GetConnString()

	m, err := migrate.New(migrationsPath, connString)
	if err != nil {
		return fmt.Errorf("failed to create migration instance :%v", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database :%v", err)
	}
	log.Println("migrated succesfully")
	return nil
}
