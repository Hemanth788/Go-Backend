package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "go.com/go-backend/db/sqlc"
)

type transferReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency, true) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency, false) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResp(errors.New("Your balance is lesser than the amount you're trying to send")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string, isFrom bool) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return false
	}

	str := fmt.Sprintf("sender account[%d] of currency: %s cannot send currency of %s", accountID, account.Currency, currency)

	if !isFrom {
		str = fmt.Sprintf("receiver account[%d] of currency: %s cannot receive currency of %s", accountID, account.Currency, currency)
	}

	if account.Currency != currency {
		err := fmt.Errorf("%s", str)
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return false
	}

	return true
}
