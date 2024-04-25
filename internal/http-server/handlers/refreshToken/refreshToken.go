package refreshToken

import (
	"Rest-shortcut/lib/api"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	refreshToken string `json:"refresh_token" validate:"required"`
}
type Response struct {
	response     *api.Response
	refreshToken string `json:"refresh_token,omitempty"`
	accessToken  string `json:"access_token,omitempty"`
}
type TokenService interface {
	ParseToken(tokenString string) (string, error)
	GenerateToken(username string) (string, error)
	GenerateRefreshToken(username string) (string, error)
}
type TokenUpdater interface {
	UpdateRefreshToken(username, refreshToken string) error
}

func New(logger *slog.Logger, tokenUpdater TokenUpdater, parser TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.refreshToken.New()"
		logger.With("op", op, "request_id", middleware.GetReqID(r.Context()))

		var request Request
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			logger.Warn("unable to decode body", "err", err)
			render.JSON(w, r, api.ErrorResponse("Fail to decode request"))
			return
		}

		err = validator.New().Struct(request)
		if err != nil {
			logger.Warn("fail to validate request", "err", err)
			render.JSON(w, r, api.ErrorResponse("Fail to validate request"))
			return
		}

		username, err := parser.ParseToken(request.refreshToken)
		if err != nil {
			logger.Warn("fail to parse refresh token", "err", err)
		}
		accessToken, err := parser.GenerateToken(username)
		if err != nil {
			logger.Warn("fail to generate token", "err", err)
			render.JSON(w, r, api.ErrorResponse("Fail to generate token"))
		}
		refreshToken, err := parser.GenerateRefreshToken(username)
		if err != nil {
			logger.Warn("fail to generate token", "err", err)
			render.JSON(w, r, api.ErrorResponse("Fail to generate token"))
		}

		err = tokenUpdater.UpdateRefreshToken(username, refreshToken)
		if err != nil {
			logger.Warn("fail to update refresh token", "err", err)
			render.JSON(w, r, api.ErrorResponse("Fail to update refresh token"))
			return
		}
		render.JSON(w, r, Response{accessToken: accessToken, refreshToken: refreshToken, response: api.SuccessResponse()})
	}
}
