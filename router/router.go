package router

import (
	"gin-web-api-ws-mqtt-frame/middleware"
	"gin-web-api-ws-mqtt-frame/service"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	api := r.Group("/api")
	api.GET("/ws", service.WsHandler)
	api.POST("/user/login", UserController.Login)

	api.Use(middleware.Jwt())
	{
		user := api.Group("/user")
		{
			user.POST("/create", UserController.Create)
		}
	}

}
