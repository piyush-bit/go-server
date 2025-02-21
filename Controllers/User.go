package controller

import (
	database "go_server/Database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		Id    int    `json:"id"`
		Email string `json:"email"`
	} `json:"data"`
}

type MyClaims struct {
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // This embeds the standard claims like exp, iat, etc.
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a hashed password with a plain text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignUp(c *gin.Context) {

	// get the email and password from the request
	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")

	// check if the email and password are empty
	if email == "" || password == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Email and Password are required",
		})
		return
	}

	// check if the email exists in the database

	isUser := database.CheckIfUserExists(email)
	if isUser {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "user already exists",
		})
		return
	}

	// hash the password
	hash, err := HashPassword(password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error hashing the password",
		})
		return
	}

	// insert the user into the database
	id, err := database.InsertUser(email, name, hash)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the user",
		})
		return
	}

	// generate a jwt token

	claim := MyClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateToken(claim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
			"id":    id,
			"email": email,
		},
	})
}

func Login(c *gin.Context) {

	// get email and password from the request
	email := c.PostForm("email")
	password := c.PostForm("password")

	// check if the email and password are empty
	if email == "" || password == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Email and Password are required",
		})
		return
	}

	// check if the email exists in the database
	user, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "user not found",
		})
		return
	}

	// check if the password is correct
	if !CheckPasswordHash(password, user.Password) {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "invalid password",
		})
		return
	}

	// generate a jwt token
	claim := MyClaims{
		Id:    user.ID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := GenerateToken(claim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
			"id":    user.ID,
			"email": email,
			"name":  user.Name,
		},
	})

}

func GetUserApps(c *gin.Context) {
	// get the user id from the request
	id := c.GetInt("id")
	// get all the apps of the user
	apps, err := database.GetAllAppsOfUser(id)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting the apps",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data":   apps,
	})
}
