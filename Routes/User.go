package routes

import (
	controller "go_server/Controllers"
	middleware "go_server/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	auth := router.Group("/api/v1")
	auth.POST("/signup", controller.SignUp)
	auth.POST("/login", controller.Login)

	// Public app routes with API key middleware
	publicApp := router.Group("/api/v1/app")
	publicApp.GET("/get/:id", controller.GetApp)

	// Protected app routes with JWT
	app := router.Group("/api/v1/app")
	app.Use(middleware.JWTAuthMiddleware())
	app.POST("/create", controller.CreateApp)
	app.GET("/list", controller.GetUserApps)
	app.PATCH("/update/:id", controller.UpdateApp)
	app.DELETE("/delete/:id", controller.DeleteApp)

	key := router.Group("/api/v1/key")
	key.GET("/public", controller.GetPublicKey)
	key.GET("/token/:id", controller.GetToken)

}
