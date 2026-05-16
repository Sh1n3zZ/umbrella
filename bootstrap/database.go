package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDatabase(config *Config) *pgxpool.Pool {
	timeout := time.Duration(config.Server.ContextTimeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dc := config.Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dc.User, dc.Password, dc.Host, dc.Port, dc.Name)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to create PostgreSQL pool: ", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping PostgreSQL: ", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return pool
}

func ClosePostgresConnection(pool *pgxpool.Pool) {
	if pool == nil {
		return
	}

	pool.Close()
	log.Println("PostgreSQL connection pool closed.")
}
