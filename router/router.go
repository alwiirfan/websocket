package router

import (
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func InitRoute(userHandler *user.Handler) {
	router = gin.Default()

	router.POST("/signup", userHandler.CreateUser)
}

func Start(addr string) error {
	return router.Run(addr)
}
