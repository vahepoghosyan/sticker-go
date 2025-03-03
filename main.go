package main

import (
    "log"
    "sticker-go/config"
    "sticker-go/routes"
)

func main() {
    config.InitDB()

    router := routes.SetupRoutes()

    log.Println("Server running on port 8080...")
    router.Run(":8080")
}
