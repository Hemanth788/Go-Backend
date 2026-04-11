package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.com/go-backend/token"
)

const (
	AUTH_HEADER_KEY  = "authorization"
	AUTH_TYPE_BEARER = "bearer"
	AUTH_PAYLOAD_KEY = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AUTH_HEADER_KEY)
		if len(authHeader) == 0 {
			err := errors.New("Auth header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid Auth header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != AUTH_TYPE_BEARER {
			err := errors.New("Auth type not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResp(err))
			return
		}

		ctx.Set(AUTH_PAYLOAD_KEY, payload)
		ctx.Next()
	}
}
