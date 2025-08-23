package routes

import (
	"material-platform/controllers"
	"material-platform/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 静态文件服务
	r.Static("/uploads", "../uploads")
	r.Static("/static", "../static")

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由（无需token）
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			// 用户相关
			user := protected.Group("/users")
			{
				user.GET("/profile", controllers.GetUserProfile)
				user.PUT("/profile", controllers.UpdateUserProfile)
			}

			// 分类管理
			categories := protected.Group("/categories")
			{
				categories.GET("/", controllers.GetCategories)
				categories.POST("/", controllers.CreateCategory)
				categories.PUT("/:id", controllers.UpdateCategory)
				categories.DELETE("/:id", controllers.DeleteCategory)
			}

			// 标签管理
			tags := protected.Group("/tags")
			{
				tags.GET("/", controllers.GetTags)
				tags.POST("/", controllers.CreateTag)
				tags.PUT("/:id", controllers.UpdateTag)
				tags.DELETE("/:id", controllers.DeleteTag)
			}

			// 文件管理
			files := protected.Group("/files")
			{
				files.GET("/", controllers.GetFiles)
				files.POST("/upload", controllers.UploadFile)
				files.GET("/:id", controllers.GetFile)
				files.PUT("/:id", controllers.UpdateFile)
				files.DELETE("/:id", controllers.DeleteFile)
				files.POST("/:id/restore", controllers.RestoreFile)
				files.GET("/:id/download", controllers.DownloadFile)
				files.GET("/:id/preview", controllers.PreviewFile)
				
				// 文件批量操作
				files.POST("/batch-delete", controllers.BatchDeleteFiles)
				files.POST("/batch-restore", controllers.BatchRestoreFiles)
			}

			// 回收站
			recycle := protected.Group("/recycle")
			{
				recycle.GET("/", controllers.GetDeletedFiles)
				recycle.DELETE("/empty", controllers.EmptyRecycleBin)
			}
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
		{
			// 用户管理
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("/", controllers.GetAllUsers)
				adminUsers.PUT("/:id", controllers.UpdateUserByAdmin)
				adminUsers.DELETE("/:id", controllers.DeleteUser)
			}

			// 系统统计
			admin.GET("/stats", controllers.GetSystemStats)
		}
	}
}