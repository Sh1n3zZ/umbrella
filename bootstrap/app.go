package bootstrap

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
)

type Application struct {
	Config *Config
	DB     *pgxpool.Pool
	Mail   *mail.Client
	Cache  *redis.Client
}

func App() Application {
	app := &Application{}
	app.Config = NewConfig()
	app.DB = NewPostgresDatabase(app.Config)
	app.Mail = NewMailClient(app.Config)
	app.Cache = NewRedisCache(app.Config)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresConnection(app.DB)
}

func (app *Application) CloseMailClient() {
	CloseMailClient(app.Mail)
}

func (app *Application) CloseRedis() {
	if app.Cache == nil {
		return
	}
	if err := app.Cache.Close(); err != nil {
		log.Println("Failed to close Redis client: ", err)
		return
	}
	log.Println("Redis client closed.")
}
