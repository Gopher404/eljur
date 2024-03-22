package httpServer

import (
	"eljur/internal/pkg/metric"
	"net/http"
)

func (h *Handler) mw(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Info("req", "URL", r.URL.String(), "Method", r.Method, "Remote", r.RemoteAddr)
		metric.HandleRequest()
		handler(w, r)
	}
}
