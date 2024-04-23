package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"net/http"
	"time"
)

const TTL = time.Hour * 24 * 30

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

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.renderTemplate(w, "", "login.html")
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
			return
		}

		login := r.Form.Get("login")
		pass := r.Form.Get("pass")
		rememberMe := r.Form.Get("remember-me") == "1"

		token, err := h.auth.Login(r.Context(), login, pass)
		if err != nil {
			h.renderTemplate(w, "Неверный логин или пароль", "login.html")
			h.l.Warn(tr.Trace(err).Error())
			return
		}

		perm, err := h.auth.GetPermission(r.Context(), login)
		if err != nil || perm < models.PermStudent {
			h.renderTemplate(w, "Недостаточно прав", "login.html")
			h.l.Warn(tr.Trace(err).Error())
			return
		}

		cookie := &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		}
		if rememberMe {
			cookie.Expires = time.Now().Add(TTL)
		}
		http.SetCookie(w, cookie)

		redirect(w, "/student/grades")
	}
}

func getToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", tr.Trace(err)
	}
	return c.Value, nil
}
