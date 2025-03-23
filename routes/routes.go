package routes

import (
	"sticker-go/config"

	"sticker-go/controllers"
	"sticker-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
    router := gin.Default()
    router.Use(config.CORSMiddleware())

    api := router.Group("/api")
    {
        api.POST("/register", controllers.RegisterUser)
        api.POST("/login", controllers.LoginUser)
        api.GET("/users", controllers.GetAllUsers)

        api.DELETE("/users", controllers.DeleteAllUsers)
        api.DELETE("/users/:id", controllers.DeleteUser)
        
        protected := api.Group("/")
        // protected.
		protected.Use(middleware.AuthMiddleware()) // Protect these routes
        protected.PUT("/profile", controllers.UpdateUserNotes)
		protected.GET("/profile", controllers.GetUserProfile)
    }

    return router
}
