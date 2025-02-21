package main

import (
	routes "go_server/Routes"

	"github.com/gin-gonic/gin"
)


func main () {
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":8080")
}
