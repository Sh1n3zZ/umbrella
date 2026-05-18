package route

import (
	"log"
	"strings"
	"time"

	"github.com/Sh1n3zZ/umbrella/api/controller"
	"github.com/Sh1n3zZ/umbrella/bootstrap"
	"github.com/Sh1n3zZ/umbrella/internal/cache"
	mail "github.com/Sh1n3zZ/umbrella/internal/mail"
	"github.com/Sh1n3zZ/umbrella/repository"
	"github.com/Sh1n3zZ/umbrella/usecase"
	gomail "github.com/wneessen/go-mail"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(
	env *bootstrap.Config,
	timeout time.Duration,
	db *pgxpool.Pool,
	mailClient *gomail.Client,
	rdb *redis.Client,
	gin *gin.Engine,
) {
	_ = timeout

	mailSender, err := mail.NewSMTPSender(mailClient, env.Mail.From)
	if err != nil {
		log.Fatal("Failed to create mail sender: ", err)
	}

	verificationStore := cache.NewEmailVerificationStore(rdb)
	verifyURLTemplate := buildVerifyURLTemplate(env.Server.PublicURL)

	base := repository.NewBaseRepository(db)
	oauthRepo := repository.NewOauthClientsRepository(base)
	oauthUC := usecase.NewOauthClientsUsecase(oauthRepo)
	oauthCtl := controller.NewOauthController(oauthUC)

	oauthUsersRepo := repository.NewOauthUsersRepository(base)
	oauthUsersUC := usecase.NewOauthUsersUsecase(oauthUsersRepo, mailSender, verificationStore, verifyURLTemplate)
	oauthUsersCtl := controller.NewOauthUsersController(oauthUsersUC)

	v1 := gin.Group("/v1")
	registerOAuthRoutes(v1, oauthCtl)
	registerOAuthUsersRoutes(v1, oauthUsersCtl)
}

// buildVerifyURLTemplate composes a fmt-style template like
// "https://app.example.com/verify?token=%s" from the configured PublicURL.
// The frontend at /verify is expected to POST the token to /v1/users/verify.
func buildVerifyURLTemplate(publicURL string) string {
	base := strings.TrimRight(strings.TrimSpace(publicURL), "/")
	if base == "" {
		return "/verify?token=%s"
	}
	return base + "/verify?token=%s"
}
