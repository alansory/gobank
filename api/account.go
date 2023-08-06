package api

import (
	"net/http"

	db "github.com/alansory/gobank/database/sqlc"
	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(http.StatusBadRequest, err, ctx)
		return
	}

	arg := db.CreateAccountParams{
		UserID:   req.UserID,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(http.StatusBadRequest, err, ctx)
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			errorResponse(http.StatusNotFound, err, ctx)
			return
		}
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, account)
}
