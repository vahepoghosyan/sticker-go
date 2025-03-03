package views

import "github.com/gin-gonic/gin"

func ResponseSuccess(message string) gin.H {
    return gin.H{"status": "success", "message": message}
}

func ResponseError(message string) gin.H {
    return gin.H{"status": "error", "message": message}
}
