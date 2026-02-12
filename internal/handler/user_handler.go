// Package handler
package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/service"
	"github.com/vcnt72/go-boilerplate/internal/utils/response"
)

type UserHandler struct {
	userService *service.UserService
}

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required"`
	Balance int64  `json:"balance" binding:"required"`
}

func (t UserHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateUserRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				response.Error(ctx, "VALIDATION_ERROR", "Salah validasi"),
			)
			return
		}

		user, err := t.userService.Create(ctx, service.CreateUserSpec{
			Balance: req.Balance,
			Name:    req.Name,
		})
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, response.Error(ctx, "UNKNOWN_ERROR", "Unknown Error"))
			return
		}

		ctx.JSON(http.StatusOK,
			response.Success(
				ctx,
				response.JSON{
					"id": user.ID,
				},
			),
		)
	}
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
