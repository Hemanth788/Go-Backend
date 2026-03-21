package api

import (
	"github.com/gin-gonic/gin"
	db "go.com/go-backend/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)

	server.router = router
	return server
}

func errorResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
