package httpServer

import (
	"eljur/internal/domain/models"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) setStudentEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/login",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermStudent, "/")
		}),
	).Methods("GET", "POST")
}
