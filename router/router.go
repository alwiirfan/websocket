package router

import (
	"server/internal/user"
	"server/internal/websocket"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func InitRoute(userHandler *user.Handler, websocketHandler *websocket.Handler) {
	router = gin.Default()

	// authentication
	router.POST("/signup", userHandler.CreateUser)
	router.POST("login", userHandler.Login)
	router.GET("/logout", userHandler.Logout)

	// websocket router
	router.POST("/websocket/createRoom", websocketHandler.CreateRoom)
}

func Start(addr string) error {
	return router.Run(addr)
}
