package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ayushrakesh/go-bank/token"

	db "github.com/ayushrakesh/go-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAccountReq

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.UserName,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, errr := server.store.GetAccount(ctx, req.ID)
	if errr != nil {
		if errr == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errr))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.UserName != account.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=3,max=5"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsReq

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.UserName,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, errr := server.store.ListAccounts(ctx, arg)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errr))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
