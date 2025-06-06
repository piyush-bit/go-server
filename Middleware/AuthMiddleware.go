package middleware

import (
	"fmt"
	controller "go_server/Controllers"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return VerifyUserToken
}
func VerifyUserToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	fmt.Println(tokenString)
	if tokenString == "" {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Authorization header is required",
		})
		c.Abort()
		return
	}
	jwtToken := strings.Split(tokenString, " ")[1]
	userClaim, err := controller.VerifyToken(jwtToken,&controller.AcessTokenClaim{})
	if err != nil {
		fmt.Println(err)
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		c.Abort()
		return
	}
	c.Set("id", userClaim.Id)
	c.Set("name", userClaim.Name)
	c.Set("email", userClaim.Email)
	c.Next()
}
