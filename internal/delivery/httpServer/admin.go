package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/internal/pkg/metric"
	"eljur/pkg/tr"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

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

	rtr.HandleFunc(url+"/metric",
		h.mw(h.handleAdminMetric),
	).Methods("GET")

	rtr.HandleFunc(url+"/logs.xlsx",
		h.mw(h.handleAdminMetricXLSX),
	).Methods("GET")

}

type adminGradesTmpData struct {
	headerTmpData
	Subjects *[4][3][]models.MinSubject `json:"subjects"`
}

func (h *Handler) handleAdminGrades(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	subjectsList, err := h.subjectService.GetAllSubjects(r.Context())
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	data := new(adminGradesTmpData)
	if err := h.SetHeaderData(r.Context(), data, "grades", login); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	data.Subjects = subjectsList

	h.renderTemplate(w, data, "admin/grades.html", "admin/header.html")
}

func (h *Handler) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	data := new(headerTmpData)
	if err := h.SetHeaderData(r.Context(), data, "users", login); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	h.renderTemplate(w, data, "admin/users.html", "admin/header.html")
}

func (h *Handler) handleAdminSubjects(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	data := new(headerTmpData)
	if err := h.SetHeaderData(r.Context(), data, "subjects", login); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	h.renderTemplate(w, data, "admin/subjects.html", "admin/header.html")
}

type metricTmpData struct {
	metric.Metric
	headerTmpData
}

func (h *Handler) handleAdminMetric(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	data := new(metricTmpData)
	m, err := metric.GetMetric()
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
	}
	data.Rps = m.Rps
	data.Logs = m.Logs
	data.RenderPerSecond = m.RenderPerSecond

	if err := h.SetHeaderData(r.Context(), data, "metric", login); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, data, "admin/metric.html", "admin/header.html")
}

func (h *Handler) handleAdminMetricXLSX(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	logs, err := metric.GetXLSXLogs()
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "logs.xlsx", time.Now(), logs)

}
