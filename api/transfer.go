package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/syucel96/simplebank/db/sqlc"
	"github.com/syucel96/simplebank/token"
	"github.com/syucel96/simplebank/util"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        string `json:"amount" binding:"required,amount"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		err := fmt.Errorf("from_account and to_account cannot have the same id: %v", req.FromAccountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	if !server.validTransfer(ctx, arg, req.Currency) {
		return
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validTransfer(ctx *gin.Context, arg db.TransferTxParams, currency string) bool {
	amount := util.ParseFloat(arg.Amount)
	if amount <= float64(0) {
		err := fmt.Errorf("transfer amount needs to be greater than 0")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	fromAccount, err := server.store.GetAccount(ctx, arg.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if fromAccount.Currency != currency {
		err = fmt.Errorf("account %d currency mismatch: expected %s, got %s", fromAccount.ID, fromAccount.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	if amount > util.ParseFloat(fromAccount.Balance) {
		err = fmt.Errorf("account %d insufficient funds: balance %s is less than requested amount %s", fromAccount.ID, fromAccount.Balance, arg.Amount)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err = fmt.Errorf("from account doesn't belong to the authenticated user %s", authPayload.Username)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return false
	}

	toAccount, err := server.store.GetAccount(ctx, arg.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if toAccount.Currency != currency {
		err = fmt.Errorf("account %d currency mismatch: expected %s, got %s", toAccount.ID, toAccount.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
