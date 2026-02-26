package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func RegisterRoutes(server *gin.Engine, handler *Handler) {
	auth := server.Group("/auth")
	auth.POST("/login", handler.Login)
}

func (h *Handler) Login(ctx *gin.Context) {
	var loginData LoginRequestDTO

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "there was an error when parsing the body"})
		return
	}

	if errs := loginData.Validate(); !errs.IsEmpty() {
		ctx.JSON(http.StatusBadRequest, errs)
		return
	}

	ctx.JSON(http.StatusOK, loginData)
}
