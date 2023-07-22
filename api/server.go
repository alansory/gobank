package api

import (
	"fmt"
	"strings"

	db "github.com/alansory/gobank/database/sqlc"
	"github.com/alansory/gobank/token"
	"github.com/alansory/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/registers", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/transfers", server.createTransfer)
	server.router = router
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
