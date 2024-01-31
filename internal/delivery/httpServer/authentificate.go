package httpServer

import "net/http"

func (h *Handler) authenticate(r *http.Request, perm int32) (login string, ok bool) {
	token, err := getToken(r)
	if err != nil {
		return "", false
	}
	login, err = h.auth.ParseToken(r.Context(), token)
	if err != nil {
		return "", false
	}

	if realPerm, err := h.auth.GetPermission(r.Context(), login); err != nil || realPerm < perm {
		return "", false
	}

	return login, true
}
