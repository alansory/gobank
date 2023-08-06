package api

import (
	"net/http"
	"regexp"
	"time"

	db "github.com/alansory/gobank/database/sqlc"
	"github.com/alansory/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type createUserRequest struct {
	Username string `json:"username"`
	Fullname string `json:"fullname" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Fullname          string    `json:"fullname"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type loginUserResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiredAt string       `json:"access_token_expired_at"`
	User                 userResponse `json:"user"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Fullname:          user.Fullname,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(http.StatusBadRequest, err, ctx)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}
	re := regexp.MustCompile(`^(.*?)@.*$`)
	username := re.ReplaceAllString(req.Email, "${1}")
	arg := db.CreateUserParams{
		Fullname:       req.Fullname,
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			errorResponse(http.StatusForbidden, err, ctx)
			return
		}
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}

	rsp := newUserResponse(user)
	successResponse(http.StatusOK, "Registration successful.", rsp, ctx)
	return
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(http.StatusBadRequest, err, ctx)
		return
	}

	user, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			errorResponse(http.StatusNotFound, err, ctx)
			return
		}
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		errorResponse(http.StatusUnauthorized, err, ctx)
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		errorResponse(http.StatusInternalServerError, err, ctx)
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
		// AccessTokenExpiredAt: accessPayload.ExpiredAt,
	}
	successResponse(http.StatusOK, "Login successful.", rsp, ctx)
	return
}
