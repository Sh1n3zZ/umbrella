package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config *Config
	DB     *pgxpool.Pool
}

func App() Application {
	app := &Application{}
	app.Config = NewConfig()
	app.DB = NewPostgresDatabase(app.Config)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresConnection(app.DB)
}
