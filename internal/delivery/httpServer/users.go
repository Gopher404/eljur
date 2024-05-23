package httpServer

import (
	"eljur/internal/domain/models"
	userService "eljur/internal/service/users"
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) setUsersEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/get_all",
		h.mw(h.handleUsersGetAll),
	).Methods("POST")

	rtr.HandleFunc(url+"/save",
		h.mw(h.handleUsersSave),
	).Methods("POST")
	rtr.HandleFunc(url+"/change_password",
		h.mw(h.handleUsersChangePassword),
	).Methods("POST")
}

func (h *Handler) handleUsersGetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	users, err := h.usersService.GetAll(r.Context())
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(users)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(out)
}

func (h *Handler) handleUsersSave(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	var in []userService.SaveUsersIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	if err := h.usersService.Save(r.Context(), in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleUsersChangePassword(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	newPassword := string(body)

	saveUsersIn := []userService.SaveUsersIn{
		{
			Action:   "update",
			Login:    login,
			Password: newPassword,
		},
	}
	if err := h.usersService.Save(r.Context(), saveUsersIn); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}
