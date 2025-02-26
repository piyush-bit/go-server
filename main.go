package main

import (
	"fmt"
	routes "go_server/Routes"
	"os"
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
	router.Static("/dist", "./dist")
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		fmt.Println(path)
		if strings.HasPrefix(path, "/api/") {
			c.JSON(404, gin.H{"error": "Not found"})
		} else if c.Request.Method == "GET" {
			c.File("./dist/index.html")
		} else {
			c.JSON(404, gin.H{"error": "Not found"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
