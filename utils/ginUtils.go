package utils

import "github.com/gin-gonic/gin"

func ErrorJson(message string) *gin.H {
	return &gin.H{"error": message}
}
