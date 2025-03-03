package controllers

import (
	"log"
	"net/http"
	"sticker-go/config"
	"sticker-go/models"
	"sticker-go/views"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)



func LoginUser(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, views.ResponseError("Invalid request"))
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, views.ResponseError("Invalid username or password"))
		return
	}

	// Log the password in both forms for debugging
	log.Println("Raw password:", request.Password)
	log.Println("Stored hashed password:", user.Password)

	// Compare the hash with the raw password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println("Error comparing password:", err)
		c.JSON(http.StatusUnauthorized, views.ResponseError("Invalid username or password"))
		return
	}

	// Generate JWT
	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

