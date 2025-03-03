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

func RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, views.ResponseError("Invalid request"))
        return
    }

    if user.Username == "" || user.Password == "" {
        c.JSON(http.StatusBadRequest, views.ResponseError("Username and Password fields are required"))
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to hash password"))
        return
    }

    log.Println("Stored password hash:", user.Password) // ✅ Debugging

    user.Password = string(hashedPassword)

    log.Println("Stored password hash:", user.Password) // ✅ Debugging

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to create user"))
        return
    }

    // Generate JWT
	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to generate token"))
		return
	}

    c.JSON(http.StatusOK, gin.H{
        "message": "User registered successfully",
        "token": token,
    })
}

func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, views.ResponseError("Unauthorized"))
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, views.ResponseError("User not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"username": user.Username,
		"notes":    user.Notes,
	})
}

func UpdateUserNotes(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, views.ResponseError("Unauthorized"))
        return
    }

    var request struct {
        Notes string `json:"notes"`
    }

    // Bind JSON request body
    log.Println(&request)
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, views.ResponseError("Invalid request format"))
        return
    }

    // Find the user by ID
    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, views.ResponseError("User not found"))
        return
    }

    // Update the notes field
    user.Notes = request.Notes
    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to update notes"))
        return
    }

    // Return success response
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "updated_notes": user.Notes,
    })
}


func GetAllUsers(c *gin.Context) {
    var users []models.User

    if err := config.DB.Select("id, username, password").Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to retrieve users"))
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "users": users})
}

func DeleteAllUsers(c *gin.Context) {
    if err := config.DB.Exec("DELETE FROM users").Error; err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to delete users"))
        return
    }

    c.JSON(http.StatusOK, views.ResponseSuccess("All users deleted successfully"))
}

func DeleteUser(c *gin.Context) {
    userID := c.Param("id") // Get user ID from the URL parameter

    if err := config.DB.Delete(&models.User{}, userID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, views.ResponseError("Failed to delete user"))
        return
    }

    c.JSON(http.StatusOK, views.ResponseSuccess("User deleted successfully"))
}