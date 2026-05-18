package main

import (
	"time"

	"github.com/Sh1n3zZ/umbrella/api/route"
	"github.com/Sh1n3zZ/umbrella/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Config

	db := app.DB
	defer app.CloseDBConnection()
	defer app.CloseMailClient()
	defer app.CloseRedis()

	timeout := time.Duration(env.Server.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, app.Mail, app.Cache, gin)

	gin.Run(env.Server.Address)
}
