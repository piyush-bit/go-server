package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	database "go_server/Database"
	"net/http"
	"strconv"
	"time"

	models "go_server/Models"

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

type AcessTokenClaim struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // This embeds the standard claims like exp, iat, etc.
}

type RefreshTokenClaim struct {
	Id                   int `json:"id"`
	jwt.RegisteredClaims     // This embeds the standard claims like exp, iat, etc.
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
	appid := c.PostForm("app_id")

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

	AccessClaim := AcessTokenClaim{
		Id:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	RefreshClaim := RefreshTokenClaim{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateToken(AccessClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	refreshToken, err := GenerateToken(RefreshClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}

	appIdInt, err := strconv.Atoi(appid)
	if err != nil {
		c.JSON(200, gin.H{
			"status": "success",
			"data": gin.H{
				"token":         token,
				"refresh_token": refreshToken,
				"id":            id,
				"email":         email,
			},
		})
		return
	}

	tokenId, err := database.InsertToken(appIdInt, token, refreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the token",
		})
		return
	}

	err = database.InsertOrUpdateSession(id, appIdInt, refreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"token_id":      tokenId,
			"id":            id,
			"email":         email,
		},
	})
}

func Login(c *gin.Context) {

	// get email and password from the request
	email := c.PostForm("email")
	password := c.PostForm("password")
	appId := c.PostForm("app_id")

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

	if user.Password == "GOOGLE" {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "User logged in with google",
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

	AccessClaim := AcessTokenClaim{
		Id:    user.ID,
		Name:  user.Name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	RefreshClaim := RefreshTokenClaim{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateToken(AccessClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	refreshToken, err := GenerateToken(RefreshClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		c.JSON(200, gin.H{
			"status": "success",
			"data": gin.H{
				"token":         token,
				"refresh_token": refreshToken,
				"id":            user.ID,
				"email":         email,
				"name":          user.Name,
			},
		})
		return
	}

	tokenId, err := database.InsertToken(appIdInt, token, refreshToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the access token",
		})
		return
	}
	err = database.InsertOrUpdateSession(user.ID, appIdInt, refreshToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the refresh token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token":         token,
			"token_id":      tokenId,
			"refresh_token": refreshToken,
			"id":            user.ID,
			"email":         email,
			"name":          user.Name,
		},
	})

}

// Google token info endpoint
const googleTokenInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

func VerifyGoogleToken(token string) (models.GoogleUser, error) {
	var user models.GoogleUser

	// Make a request to Google's token verification endpoint
	req, err := http.NewRequest("GET", googleTokenInfoURL, nil)
	if err != nil {
		return user, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return user, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return user, fmt.Errorf("token verification failed with status %d: %s", resp.StatusCode, string(body[:n]))
		
	}

	// Decode the JSON response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Optional: You can add further validation here (e.g., check audience/client_id)

	return user, nil
}

func ContinueWithGoogle(c *gin.Context) {
	GoogleToken := c.PostForm("google_token")
	if GoogleToken == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Google token is required",
		})
		return
	}
	// verify the google token
	googleUser, err := VerifyGoogleToken(GoogleToken)
	if err != nil {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid google token",
		})
		return
	}
	// check if the user exists in the database
	user, err := database.GetUserByEmail(googleUser.Email)
	if err != nil {
		// if the user does not exist, create a new user
		id, err := database.InsertUser(googleUser.Email, googleUser.Name, "GOOGLE")
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Error inserting the user",
			})
			return
		}
		user = models.User{
			ID:    id,
			Name:  googleUser.Name,
			Email: googleUser.Email,
		}
	}
	// generate a jwt token
	AccessClaim := AcessTokenClaim{
		Id:    user.ID,
		Name:  googleUser.Name,
		Email: googleUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	RefreshClaim := RefreshTokenClaim{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GenerateToken(AccessClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	refreshToken, err := GenerateToken(RefreshClaim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	appId := c.PostForm("app_id")
	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		c.JSON(200, gin.H{
			"status": "success",
			"data": gin.H{
				"token":         token,
				"refresh_token": refreshToken,
				"id":            user.ID,
				"email":         user.Email,
				"name":          user.Name,
			},
		})
		return
	}
	tokenId, err := database.InsertToken(appIdInt, token, refreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the token",
		})
		return
	}
	err = database.InsertOrUpdateSession(user.ID, appIdInt, refreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error inserting the token",
		})
		return
	}
	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"token_id":      tokenId,
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
		},
	})
}

func ChangePassword(c *gin.Context) {
	// get the token from the request
	id := c.PostForm("id")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	// check if the token is valid
	user, err := database.GetUserById(id)
	if err != nil {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		return
	}
	// check if the old password is correct
	if !CheckPasswordHash(oldPassword, user.Password) {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid password",
		})
		return
	}

	// hash the new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error hashing the password",
		})
		return
	}

	// update the password in the database
	err = database.UpdatePassword(id, hashedPassword)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating the password",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Password changed successfully",
	})
}

func Refresh(c *gin.Context) {
	// get the token from the request
	token := c.PostForm("token")
	id := c.PostForm("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid id",
		})
		return
	}
	// check if the token is empty
	if token == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
		return
	}

	// check if the token is valid
	claims, err := VerifyToken(token, &RefreshTokenClaim{})
	if err != nil {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		return
	}

	refreshToken, user, err := database.GetRefreshToken(idInt, claims.Id)

	if err != nil {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		return
	}

	if refreshToken != token {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Invalid token",
		})
		return
	}

	// generate a new token
	claim := RefreshTokenClaim{
		Id: claims.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newToken, err := GenerateToken(claim)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}
	newAccessToken, err := GenerateToken(AcessTokenClaim{
		Id:    claims.Id,
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error generating the token",
		})
		return
	}

	// update the token in the database
	err = database.UpdateRefreshToken(idInt, claims.Id, newToken)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error updating the token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token":         newAccessToken,
			"refresh_token": newToken,
			"id":            claims.Id,
			"email":         user.Email,
			"name":          user.Name,
		},
	})
}

func Logout(c *gin.Context) {
	// get the user id from the request
	id := c.GetInt("id")
	appId := c.PostForm("app_id")

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Invalid app_id",
		})
		return
	}

	// delete the token from the database
	err = database.DeleteSession(id, appIdInt)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error deleting the token",
		})
		return
	}

	// send the response
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Logout successful",
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

func Home(c *gin.Context) {
	id := c.GetInt("id")
	apps, err := database.GetAllAppsOfUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": "Apps not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Error getting the apps",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"apps": apps,
			"user": gin.H{
				"id":    id,
				"name":  c.GetString("name"),
				"email": c.GetString("email"),
			},
		},
	})
}
