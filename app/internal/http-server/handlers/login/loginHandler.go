package login

import (
	"Rest-shortcut/lib/api"
	sl "Rest-shortcut/lib/logger"
	"Rest-shortcut/lib/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type Response struct {
	Response     api.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type UserGetter interface {
	GetUser(login, password string) (models.User, error)
}

type TokenGenerator interface {
	GenerateToken(username string) (string, error)
	GenerateRefreshToken(username string) (string, error)
}

func New(logger *slog.Logger, storage UserGetter, auth TokenGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login.New"
		logger.With("op", op,
			"request_id", middleware.GetReqID(r.Context()))

		var request Request
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			logger.Error("Fail to decode request", sl.Err(err))
			render.JSON(w, r, api.ErrorResponse("fail to decode request"))
			return
		}

		err = validator.New().Struct(request)
		if err != nil {
			logger.Error("Fail to validate request", sl.Err(err))
			render.JSON(w, r, api.ErrorResponse("fail to validate request"))
			return
		}

		user, err := storage.GetUser(request.Login, request.Password)
		if err != nil || user.Password != request.Password {
			logger.Error("Fail to login", sl.Err(err))
			render.JSON(w, r, api.ErrorResponse("Login or password incorrect"))
			return
		}
		accessToken, err := auth.GenerateToken(user.Login)
		refreshToken, err := auth.GenerateRefreshToken(user.Login)

		response := Response{Response: *api.SuccessResponse(), AccessToken: accessToken, RefreshToken: refreshToken}
		logger.Info("Success login", user.Login)
		render.JSON(w, r, response)
	}
}
