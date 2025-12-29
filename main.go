package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", HealthHandler)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
