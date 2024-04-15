package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"video_convert_tool/internal/http-server/handlers/convert"
)

func InitRouter(log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/convert_video", func(r chi.Router) {
		r.Post("/", convert.ConvertVideo(log))
	})

	return router
}
