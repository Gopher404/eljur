package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) setAdminEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/grades",
		h.logHandle(h.handleAdminGrades),
	).Methods("GET")

	rtr.HandleFunc(url+"/users",
		h.logHandle(h.handleAdminUsers),
	).Methods("GET")

	rtr.HandleFunc(url+"/login",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermAdmin, url+"/grades")
		}),
	).Methods("GET", "POST")

}

type adminGradesTmpData struct {
	Subjects []models.Subject `json:"subjects"`
}

func (h *Handler) handleAdminGrades(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/admin/login")
		return
	}
	subjectsList, err := h.subjectService.GetAllSubjects()
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	h.renderTemplate(w, "admin/grades.html", adminGradesTmpData{
		Subjects: subjectsList,
	})
}

func (h *Handler) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/admin/login")
		return
	}
	h.renderTemplate(w, "admin/users.html", nil)
}
