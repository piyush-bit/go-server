package routes

import (
	controller "go_server/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	public := router.Group("/api/v1")
	public.POST("/signup", controller.SignUp)
	public.POST("/login", controller.Login)

}
