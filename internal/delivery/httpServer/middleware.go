package httpServer

import (
	"context"
	"eljur/internal/pkg/metric"
	"net/http"
)

func (h *Handler) mw(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()

		r = r.WithContext(ctx)
		h.l.Info("req", "URL", r.URL.String(), "Method", r.Method, "Remote", r.RemoteAddr)
		metric.HandleRequest()

		handler(w, r)
	}
}
