package controller

import (
	"database/sql"
	database "go_server/Database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) {
	id := c.Param("id")
	tokenId,err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid token ID",
		})
		return
	}
	// get the token from the database
	userToken , err := database.GetTokenById(tokenId)
	if err != nil {
		if(err == sql.ErrNoRows){
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "Token not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting the token",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data" : userToken,
	})

}
