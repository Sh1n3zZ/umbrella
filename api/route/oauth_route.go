package route

import (
	"github.com/Sh1n3zZ/umbrella/api/controller"
	"github.com/gin-gonic/gin"
)

func registerOAuthRoutes(v1 *gin.RouterGroup, oauthCtl *controller.OauthController) {
	oauth := v1.Group("/oauth")
	oauth.GET("/authorize", oauthCtl.Authorization)
}
