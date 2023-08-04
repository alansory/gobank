package api

import (
	"fmt"
	"io"
	"net/http"
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

	report := gin.H{
		"error": gin.H{
			"status_code": http.StatusInternalServerError,
			"message":     err.Error(),
		},
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		report["error"].(gin.H)["status_code"] = http.StatusUnprocessableEntity
		report["error"].(gin.H)["message"] = "422 Unprocessable Entity"
		errors := make(map[string]string)

		for _, validationErr := range validationErrors {
			switch validationErr.Tag() {
			case "required":
				errors[strings.ToLower(validationErr.Field())] = fmt.Sprintf("The %s field is required", strings.ToLower(validationErr.Field()))
			case "email":
				errors[strings.ToLower(validationErr.Field())] = fmt.Sprintf("The %s field is not a valid email", strings.ToLower(validationErr.Field()))
			case "gte":
				errors[strings.ToLower(validationErr.Field())] = fmt.Sprintf("The %s field value must be greater than %s", strings.ToLower(validationErr.Field()), validationErr.Param())
			case "lte":
				errors[strings.ToLower(validationErr.Field())] = fmt.Sprintf("The %s field value must be lower than %s", strings.ToLower(validationErr.Field()), validationErr.Param())
			}
		}
		report["error"].(gin.H)["errors"] = errors
	}

	if err == io.EOF {
		report["error"].(gin.H)["status_code"] = http.StatusBadRequest
		report["error"].(gin.H)["message"] = "JSON data is missing or malformed"
	}

	return report
}
