package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(data any) map[string]any {
	return gin.H{
		"code":    200,
		"data":    data,
		"message": "success",
	}
}

func ErrorResponse(data any, error string) map[string]any {
	return gin.H{
		"code":    500,
		"data":    data,
		"message": error,
	}
}

//func Success_for_page() map[string]any {
//	return gin.H{
//		"code":    500,
//		"data":    data,
//		"message": error,
//	}
//}
