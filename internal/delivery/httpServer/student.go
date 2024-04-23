package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) setStudentEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/grades",
		h.mw(h.handleStudentGrades),
	).Methods("GET")
}

type StudentGradesTmpData struct {
	headerTmpData
	Subjects *[4][3][]models.MinSubject `json:"subjects"`
}

func (h *Handler) handleStudentGrades(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login")
		return
	}
	data := new(StudentGradesTmpData)
	if err := h.SetHeaderData(r.Context(), data, "grades", login); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	subjects, err := h.subjectService.GetAllSubjects(r.Context())
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	data.Subjects = subjects
	h.renderTemplate(w, data, "/student/grades.html", "/student/header.html")
}
