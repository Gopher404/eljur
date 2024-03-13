package httpServer

import (
	"eljur/internal/domain/models"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var studentUrl string

func (h *Handler) setStudentEndpoints(rtr *mux.Router, url string) {
	studentUrl = url

	rtr.HandleFunc(url+"/login",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermStudent, url+"/grades")
		}),
	).Methods("GET", "POST")

	rtr.HandleFunc(url+"/grades",
		h.logHandle(h.handleStudentGrades),
	).Methods("GET")
}

func (h *Handler) handleStudentGrades(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermStudent); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, studentUrl+"/login")
		return
	}
	h.renderTemplate(w, "/student/grades.html", nil)
}
