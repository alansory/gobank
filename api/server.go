package api

import (
	"strings"

	db "github.com/alansory/gobank/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/registers", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error, ctx *gin.Context) gin.H {
	if ctx.Request.Method == "POST" {
		if ctx.Request.ContentLength == 0 || !strings.Contains(ctx.ContentType(), "application/json") {
			return gin.H{"error": "Invalid or empty JSON data"}
		}
	}
	return gin.H{"error": err.Error()}
}
