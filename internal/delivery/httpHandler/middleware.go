package httpHandler

import (
	"bytes"
	"context"
	"eljur/internal/pkg/metric"
	"io"
	"net/http"
)

func (h *Handler) mw(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		defer cancel()

		r = r.WithContext(ctx)

		body, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		if err != nil {
			body = []byte{}
		}

		metric.HandleRequest()
		h.l.Info("req", "URL", r.URL.String(), "Method", r.Method, "Remote", r.RemoteAddr, "RPS", metric.GetRPS(), "body", string(body))

		handler(w, r)
	}
}
