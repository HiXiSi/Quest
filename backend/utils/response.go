package utils

import (
	"material-platform/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.ApiResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusBadRequest, models.ApiResponse{
		Code:    code,
		Message: message,
	})
}

// ServerErrorResponse 服务器错误响应
func ServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, models.ApiResponse{
		Code:    500,
		Message: message,
	})
}

// UnauthorizedResponse 未授权响应
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, models.ApiResponse{
		Code:    401,
		Message: message,
	})
}

// ForbiddenResponse 禁止访问响应
func ForbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, models.ApiResponse{
		Code:    403,
		Message: message,
	})
}

// NotFoundResponse 未找到响应
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, models.ApiResponse{
		Code:    404,
		Message: message,
	})
}

// PageResponse 分页响应
func PageResponse(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Code:    0,
		Message: "success",
		Data: models.PageResponse{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Pages:    pages,
		},
	})
}