package main

import (
	"github/afiffaizun/matric-api/internal/handler"
	"github/afiffaizun/matric-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.PrometheusMetrics())

	r.GET("/", handler.HealthCheck)
	r.GET("/health", handler.HealthCheck)
	r.GET("/api/matrix", handler.GetMatrix)
	r.POST("/api/matrix", handler.CreateMatrix)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8080")
}
