package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var studentUrl string

func (h *Handler) setStudentEndpoints(rtr *mux.Router, url string) {
	studentUrl = url

	rtr.HandleFunc(url+"/login",
		h.mw(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermStudent, url+"/grades")
		}),
	).Methods("GET", "POST")

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
		redirect(w, studentUrl+"/login")
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
