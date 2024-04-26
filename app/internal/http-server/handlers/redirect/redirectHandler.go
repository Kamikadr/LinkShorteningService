package redirect

import (
	"Rest-shortcut/lib/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type UrlGetter interface {
	GetUrl(user, longUrl string) (string, error)
}

func New(logger *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.With("request_id", middleware.GetReqID(r.Context()))
		user := r.Context().Value("user").(string)
		shortUrl := chi.URLParam(r, "text")
		if shortUrl == "" {
			logger.Warn("shortUrl is empty")
			render.JSON(w, r, api.ErrorResponse("shortUrl is empty"))
			return
		}
		longUrl, err := urlGetter.GetUrl(user, shortUrl)
		if err != nil {
			logger.Info("url not found", shortUrl)
			render.JSON(w, r, api.ErrorResponse("url not found"))
		}
		logger.Info("Find url", shortUrl, longUrl)
		http.Redirect(w, r, longUrl, http.StatusFound)
	})

}
