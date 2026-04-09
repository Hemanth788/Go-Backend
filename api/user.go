package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "go.com/go-backend/db/sqlc"
	"go.com/go-backend/util"
)

type createUserReq struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResp struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserReq struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserResp struct {
	AccessToken string `json:"access_token"`
	User userResp `json:"user"`
}

func newUserResponse(user db.User) userResp {
	return userResp{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		HashedPassword: hashedPassword,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResp(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, newUserResponse(user))
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResp(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	response := loginUserResp{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, response)
}