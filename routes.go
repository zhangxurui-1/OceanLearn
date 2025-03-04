package main

import (
	"github.com/gin-gonic/gin"
	"oceanlearn/controller"
	"oceanlearn/middleware"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	// 文章分类
	categoryRoutes := r.Group("/categories")
	categoryController := controller.NewCategoryController()

	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update)
	categoryRoutes.GET("/:id", categoryController.Show)
	categoryRoutes.DELETE("/:id", categoryController.Delete)

	// 文章
	postRoutes := r.Group("/posts")
	postController := controller.NewPostController()

	postRoutes.POST("", PostController.Create)
	postRoutes.PUT("/:id", PostController.Update)
	postRoutes.GET("/:id", PostController.Show)
	postRoutes.DELETE("/:id", PostController.Delete)

	return r
}
