package httpServer

import (
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

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request, perm int32, outUrl string) {
	if r.Method == "GET" {
		h.renderTemplate(w, "login.html", Message{})
	} else {
		if err := r.ParseForm(); err != nil {
			h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
			return
		}

		login := r.Form.Get("login")
		pass := r.Form.Get("pass")

		token, err := h.auth.Login(r.Context(), login, pass)
		if err != nil {
			h.renderTemplate(w, "login.html", Message{Mess: "Неверный логин или пароль"})
			h.l.Warn(tr.Trace(err).Error())
			return
		}

		realPerm, err := h.auth.GetPermission(r.Context(), login)
		if err != nil || realPerm < perm {
			h.renderTemplate(w, "login.html", Message{Mess: "Недостаточно прав"})
			h.l.Warn(tr.Trace(err).Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(TTL),
		})

		redirect(w, outUrl)
	}
}

func getToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", tr.Trace(err)
	}
	return c.Value, nil
}
