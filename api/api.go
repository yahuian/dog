package api

import (
	"dog/config"
	"dog/docs"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

func Register(e *gin.Engine) {
	api := e.Group("/api/v1")

	if config.Get().Server.Mode == config.ServerModeDebug {
		// API 文档 /api/docs/index.html
		docs.SwaggerInfo.BasePath = "/api/v1"
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	controllers := []func(*gin.RouterGroup){
		// 各模块路由
	}
	for _, register := range controllers {
		register(api)
	}
}
