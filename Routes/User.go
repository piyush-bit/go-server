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
	auth.POST("/refresh", controller.Refresh)
	auth.POST("/change-password",controller.ChangePassword)
	auth.POST("/forget-password", controller.InitiateForgetPassword)
	auth.POST("/reset-password", controller.CompleteForgetPassword)
	auth.POST("/google-login", controller.ContinueWithGoogle)


	// Protected user routes with JWT
	auth.Use(middleware.JWTAuthMiddleware())
	auth.POST("/logout", controller.Logout)

	// Public app routes with API key middleware
	publicApp := router.Group("/api/v1/app")
	publicApp.GET("/get/:id", controller.GetApp)

	// Protected app routes with JWT
	app := router.Group("/api/v1/app")
	app.Use(middleware.JWTAuthMiddleware())
	app.GET("/", controller.Home)
	app.POST("/create", controller.CreateApp)
	app.GET("/list", controller.GetUserApps)
	app.PATCH("/:id", controller.UpdateApp)
	app.DELETE("/:id", controller.DeleteApp)

	key := router.Group("/api/v1/key")
	key.GET("/public", controller.GetPublicKey)
	key.GET("/token/:id", controller.GetToken)

}
