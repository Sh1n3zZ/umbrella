package route

import (
	"time"

	"github.com/Sh1n3zZ/umbrella/api/controller"
	"github.com/Sh1n3zZ/umbrella/bootstrap"
	"github.com/Sh1n3zZ/umbrella/repository"
	"github.com/Sh1n3zZ/umbrella/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(env *bootstrap.Config, timeout time.Duration, db *pgxpool.Pool, gin *gin.Engine) {
	_ = env
	_ = timeout

	base := repository.NewBaseRepository(db)
	oauthRepo := repository.NewOauthClientsRepository(base)
	oauthUC := usecase.NewOauthClientsUsecase(oauthRepo)
	oauthCtl := controller.NewOauthController(oauthUC)

	v1 := gin.Group("/v1")
	registerOAuthRoutes(v1, oauthCtl)
}
