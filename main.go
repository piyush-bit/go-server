package main

import (
	routes "go_server/Routes"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	routes.SetupRoutes(router)
	router.Static("/assets", "./dist/assets")
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api/") {
			c.File("./dist/index.html")
		} else {
			c.JSON(404, gin.H{"message": "Page not found"})
		}
	})
	router.Run(":8080")
}
