package httpServer

import (
	"eljur/pkg/tr"
	"net/http"
)

func (h *Handler) authenticate(r *http.Request, perm int32) (login string, ok bool) {
	token, err := getToken(r)
	if err != nil {
		h.l.Warn(tr.Trace(err).Error())
		return "", false
	}
	login, err = h.auth.ParseToken(r.Context(), token)
	if err != nil {
		h.l.Warn(tr.Trace(err).Error())
		return "", false
	}

	realPerm, err := h.auth.GetPermission(r.Context(), login)
	if err != nil {
		h.l.Warn(tr.Trace(err).Error())
		return "", false
	}
	if realPerm < perm {
		return "", false
	}

	return login, true
}
