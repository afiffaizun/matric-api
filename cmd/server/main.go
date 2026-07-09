package main

import (
	"github/afiffaizun/matric-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", handler.HealthCheck)
	r.GET("/health", handler.HealthCheck)
	r.GET("/api/matrix", handler.GetMatrix)
	r.POST("/api/matrix", handler.CreateMatrix)

	r.Run(":8080")
}
