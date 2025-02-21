package controller

import (
	"database/sql"
	database "go_server/Database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateApp(c *gin.Context) {
	name := c.PostForm("name")
	callback_url := c.PostForm("callback_url")

	// check if the name and callback_url are empty
	if name == "" || callback_url == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Name and Callback URL are required",
		})
		return
	}

	// get the user id from the context
	id, _ := c.Get("id")

	// insert the app into the database
	appId, err := database.InsertApp(name, callback_url, id.(int))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the app",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"id":           appId,
			"name":         name,
			"callback_url": callback_url,
		},
	})

}

func GetApp(c *gin.Context) {
	id := c.Param("id")
	// parse the id into an integer
	appId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid app ID",
		})
		return
	}

	// get the user id from the context
	userId, _ := c.Get("id")

	// get the app from the database
	app, err := database.GetAppById(appId)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting the app",
		})
		return
	}

	if(userId != app.UserId){
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   app,
	})
}

func UpdateApp(c *gin.Context){
	id := c.Param("id")
	// parse the id into an integer
	appId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid app ID",
		})
		return
	}

	// get the user id from the context
	userId, _ := c.Get("id")
	userIdInt := userId.(int)

	name := c.PostForm("name")
	callback_url := c.PostForm("callback_url")

	// check if the name and callback_url are empty
	if name == "" || callback_url == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Name and Callback URL are required",
		})
		return
	}

	err = database.UpdateApp(appId,userIdInt, name, callback_url)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "App not found",
			})
		}
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating the app",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"id":           appId,
			"name":         name,
			"callback_url": callback_url,
		},
	})
}

func DeleteApp(c *gin.Context){
	id := c.Param("id")
	// parse the id into an integer
	appId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid app ID",
		})
		return
	}

	// get the user id from the context
	userId, _ := c.Get("id")
	userIdInt := userId.(int)

	err = database.DeleteApp(appId,userIdInt)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "App not found",
			})
		}
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting the app",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"message": "App deleted successfully",
	})
}
