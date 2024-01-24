package httpServer

import (
	"fmt"
	"net/http"
)

func (h *Handler) logHandle(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.l.Info(fmt.Sprintf("%s %s, remote: %s", r.URL, r.Method, r.RemoteAddr))
		handler(w, r)
	}
}
