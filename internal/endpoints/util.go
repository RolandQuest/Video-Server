package endpoints

import (
	"log/slog"
	"net/http"
	"github.com/a-h/templ"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func WrapError(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	if err := component.Render(r.Context(), w); err != nil {
		return err
	}
	return nil
}