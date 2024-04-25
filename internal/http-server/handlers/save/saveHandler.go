package save

import (
	"Rest-shortcut/lib/api"
	sl "Rest-shortcut/lib/logger"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	OldUrl string `json:"old_url" validate:"required,url"`
	NewUrl string `json:"new_url" validate:"required"`
}
type UrlSaver interface {
	SaveUrl(user, oldUrl, shortUrl string) error
}

func New(logger *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.With(slog.String("request_id", middleware.GetReqID(request.Context())))
		user := request.Context().Value("user").(string)
		var req *Request
		err := render.DecodeJSON(request.Body, &req)
		if errors.Is(err, io.EOF) {
			logger.Warn("request body is empty")
			render.JSON(writer, request, api.ErrorResponse("empty request"))
			return
		}
		if err != nil {
			logger.Error("Fail to deserialize request", sl.Err(err))
			render.JSON(writer, request, api.ErrorResponse("Error deserializing request"))
			return
		}
		if err := validator.New().Struct(*req); err != nil {
			var validatorError validator.ValidationErrors
			errors.As(err, &validatorError)
			logger.Error("Validation error", sl.Err(validatorError))
			render.JSON(writer, request, api.ErrorResponse(validatorError.Error()))
		}

		if err := urlSaver.SaveUrl(user, req.OldUrl, req.NewUrl); err != nil {
			logger.Error("Fail to save url", sl.Err(err))
			render.JSON(writer, request, api.ErrorResponse("Fail to save url"))
			return
		}
		logger.Info("Url saved")
		render.JSON(writer, request, api.SuccessResponse())
	}
}
