package controller

import (
	database "go_server/Database"
	services "go_server/Services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ForgetPasswordClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func InitiateForgetPassword(c *gin.Context) {

	// get user by email to vertify if user exists
	email := c.PostForm("email")
	if email == "" {
		c.JSON(400, gin.H{
			"message": "email is required",
		})
	}

	user, err := database.GetUserByEmail(email)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
	}

	claim := ForgetPasswordClaim{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "go_server",
			Subject:   "forget_password",
		},
	}

	token, err := GenerateToken(claim)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "error generating token",
		})
	}
	//insert token to db
	err = database.InsertForgetPassword(email, token, time.Now().Add(time.Hour))
	if err!= nil {
		c.JSON(400, gin.H{
			"message": "error inserting token",
		})
	}
	// get the backend url
	link := c.Request.Host + "/complete-forget-password?email=" + user.Email + "&token=" + token

	err = services.SendForgetPasswordEmail(user.Email, link)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "error sending email",
		})
	}
	c.JSON(200, gin.H{
		"message": "email sent",
	})
}

func CompleteForgetPassword(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")
	newPassword := c.PostForm("password")
	if email == "" || token == "" || newPassword == "" {
		c.JSON(400, gin.H{
			"message": "email ,token and password are required",
		})
	}
	claim := &ForgetPasswordClaim{}
	claim, err := VerifyToken(token, claim)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error verifying token",
		})
	}
	if claim.Email != email {
		c.JSON(400, gin.H{
			"message": "email does not match",
		})
	}

	// hash the password
	hashPassword, err := HashPassword(newPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error hashing password",
		})
	}

	// update the password
	err = database.UpdatePassword(email, hashPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error updating password",
		})
	}
	c.JSON(200, gin.H{
		"message": "password updated",
	})
}
