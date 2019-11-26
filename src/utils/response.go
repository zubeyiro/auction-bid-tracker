package utils

import (
	"github.com/gin-gonic/gin"
)

// NewHTTPResponse is global HTTP response function, status represents if process completed successfully
func NewHTTPResponse(status bool, data interface{}) map[string]interface{} {
	return gin.H{"status": status, "data": data}
}
