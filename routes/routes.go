package routes

import (
	"devso-backend/controllers"
	"devso-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 使用中间件
	router.Use(middlewares.Logger())

	// 路由组
	v1 := router.Group("/api/v1")
	{
		v1.GET("/search", controllers.Search)
	}

	return router
}
