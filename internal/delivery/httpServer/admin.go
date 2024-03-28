package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/internal/pkg/metric"
	"eljur/pkg/tr"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var adminUrl string

func (h *Handler) setAdminEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/grades",
		h.mw(h.handleAdminGrades),
	).Methods("GET")

	rtr.HandleFunc(url+"/users",
		h.mw(h.handleAdminUsers),
	).Methods("GET")

	rtr.HandleFunc(url+"/subjects",
		h.mw(h.handleAdminSubjects),
	).Methods("GET")

	rtr.HandleFunc(url+"/login",
		h.mw(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermAdmin, url+"/grades")
		}),
	).Methods("GET", "POST")
	rtr.HandleFunc(url+"/metric",
		h.mw(h.handleAdminMetric),
	).Methods("GET")

}

type adminGradesTmpData struct {
	Subjects *[4][3][]models.MinSubject `json:"subjects"`
}

func (h *Handler) handleAdminGrades(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, adminUrl+"/login")
		return
	}
	subjectsList, err := h.subjectService.GetAllSubjects(r.Context())
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
		redirect(w, adminUrl+"/login")
		return
	}
	h.renderTemplate(w, "admin/users.html", nil)
}

func (h *Handler) handleAdminSubjects(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, adminUrl+"/login")
		return
	}
	h.renderTemplate(w, "admin/subjects.html", nil)
}

func (h *Handler) handleAdminMetric(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, adminUrl+"/login")
		return
	}
	m, err := metric.GetMetric()
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
	}
	h.renderTemplate(w, "admin/metric.html", m)
}
