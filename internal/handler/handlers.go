package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "matric-api is running",
	})
}

func GetMatrix(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status: "ok",
		Data:   "TODO: implement matrix retrieval",
	})
}

func CreateMatrix(c *gin.Context) {
	c.JSON(http.StatusCreated, Response{
		Status:  "ok",
		Message: "TODO: implement matrix creation",
	})
}
