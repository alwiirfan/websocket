package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (handler *Handler) CreateUser(ctx *gin.Context) {
	var request CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := handler.Service.CreateUser(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (handler *Handler) Login(ctx *gin.Context) {
	var user LoginUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, err := handler.Service.Login(ctx.Request.Context(), &user)
	if err != nil {
		fmt.Println("Error in login:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("jwt", usr.accessToken, 3600, "/", "localhost", false, true)

	respon := &LoginUserResponse{
		Username: usr.Username,
		ID:       usr.ID,
	}

	ctx.JSON(http.StatusOK, respon)
}

func (handler *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
