package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ayushrakesh/go-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	var req createTransferReq

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validateAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) validateAccount(ctx *gin.Context, id int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		errr := fmt.Errorf("account %d mismatch found currency %s vs %s", id, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(errr))
		return false
	}
	return true
}
