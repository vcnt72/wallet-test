package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/domain"
	"github.com/vcnt72/go-boilerplate/internal/service"
	"github.com/vcnt72/go-boilerplate/internal/utils/response"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func (w WalletHandler) GetBalance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.GetHeader("X-User-ID")

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(ctx, "INVALID_USER_ID", "user_id must be a positive integer"))
			return
		}

		wallet, err := w.walletService.GetByUserID(ctx, userID)
		if err != nil {

			if errors.Is(err, domain.ErrWalletNotFound) {
				ctx.JSON(http.StatusNotFound, response.Error(ctx, "DATA_NOT_FOUND", "User not found"))
				return
			}

			ctx.JSON(http.StatusInternalServerError, response.Error(ctx, "UNKNOWN_ERROR", "Unknown error"))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(ctx, response.JSON{
			"balance": wallet.Balance,
		}))
	}
}

type WithdrawRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

func (w WalletHandler) Withdraw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req WithdrawRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				response.Error(ctx, "VALIDATION_ERROR", "Salah validasi"),
			)
			return
		}

		userIDStr := ctx.GetHeader("X-User-ID")

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(ctx, "INVALID_USER_ID", "user_id must be a positive integer"))
			return
		}

		idempotencyKey := ctx.GetHeader("X-Idempotency-Key")

		if idempotencyKey == "" {
			ctx.JSON(http.StatusBadRequest, response.Error(ctx, "INVALID_IDEMPOTENCY_KEY", "idempotency should exist"))
			return

		}

		withdrawalRes, err := w.walletService.Withdraw(ctx, service.WithdrawWalletSpec{
			UserID:         userID,
			IdempotencyKey: idempotencyKey,
			Amount:         req.Amount,
		})
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrInvalidAmount):
				ctx.JSON(http.StatusBadRequest,
					response.Error(ctx, "INVALID_AMOUNT", "amount must be greater than 0"))
				return

			case errors.Is(err, domain.ErrWalletNotFound):
				ctx.JSON(http.StatusNotFound,
					response.Error(ctx, "WALLET_NOT_FOUND", "wallet not found"))
				return

			case errors.Is(err, domain.ErrInsufficientFund):
				ctx.JSON(http.StatusConflict,
					response.Error(ctx, "INSUFFICIENT_FUNDS", "insufficient balance"))
				return

			case errors.Is(err, domain.ErrIdempotencyKeyReused):
				ctx.JSON(http.StatusConflict,
					response.Error(ctx, "IDEMPOTENCY_KEY_REUSED", "idempotency key reused with different request"))
				return

			case errors.Is(err, domain.ErrRequestInProgress):
				ctx.JSON(http.StatusConflict,
					response.Error(ctx, "REQUEST_IN_PROGRESS", "request is being processed, please retry"))
				return

			case errors.Is(err, domain.ErrWithdrawFailed):
				ctx.JSON(http.StatusInternalServerError,
					response.Error(ctx, "WITHDRAW_FAILED", "withdraw failed"))
				return

			default:
				ctx.JSON(http.StatusInternalServerError,
					response.Error(ctx, "UNKNOWN_ERROR", "internal server error"))
				return
			}
		}

		ctx.JSON(http.StatusOK, response.Success(ctx, response.JSON{
			"balance": withdrawalRes.Balance,
			"amount":  withdrawalRes.Amount,
			"userId":  withdrawalRes.UserID,
		}))
	}
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}
