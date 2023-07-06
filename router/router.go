package router

import (
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func InitRoute(userHandler *user.Handler) {
	router = gin.Default()

	router.POST("/signup", userHandler.CreateUser)
	router.POST("login", userHandler.Login)
	router.GET("/logout", userHandler.Logout)
}

func Start(addr string) error {
	return router.Run(addr)
}
