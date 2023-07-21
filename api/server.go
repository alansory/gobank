package api

import (
	"strings"

	db "github.com/alansory/gobank/database/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error, ctx *gin.Context) gin.H {
	if ctx.Request.ContentLength == 0 || !strings.Contains(ctx.ContentType(), "application/json") {
		return gin.H{"error": "Invalid or empty JSON data"}
	}
	return gin.H{"error": err.Error()}
}
