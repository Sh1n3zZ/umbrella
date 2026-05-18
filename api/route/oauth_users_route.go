package route

import (
	"github.com/Sh1n3zZ/umbrella/api/controller"
	"github.com/gin-gonic/gin"
)

func registerOAuthUsersRoutes(v1 *gin.RouterGroup, oauthUsersCtl *controller.OauthUsersController) {
	users := v1.Group("/users")
	users.POST("/register", oauthUsersCtl.Register)
	users.POST("/verify", oauthUsersCtl.Verify)
}
