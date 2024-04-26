package register

import (
	"Rest-shortcut/lib/api"
	sl "Rest-shortcut/lib/logger"
	"Rest-shortcut/storage/postrges"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func New(logger *slog.Logger, storage *postrges.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register.New"
		logger.With(slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Fail to decode request", sl.Err(err))
			render.JSON(w, r, api.ErrorResponse("Fail to decode request"))
			return
		}
		err = validator.New().Struct(req)
		if err != nil {
			var validationError validator.ValidationErrors
			errors.As(err, validationError)
			logger.Error("Fail to validate request", sl.Err(validationError))
			render.JSON(w, r, api.ErrorResponse(validationError.Error()))
			return
		}

		err = storage.AddUser(req.Login, req.Password)
		if err != nil {
			logger.Error("Fail to add user", sl.Err(err))
			render.JSON(w, r, api.ErrorResponse("Fail to add user"))
			return
		}
		logger.Info("Successfully registered user")
		render.JSON(w, r, api.SuccessResponse())
		return
	}
}
