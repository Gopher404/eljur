package httpHandler

import (
	"eljur/internal/domain/models"
	"eljur/internal/service/grades"
	"eljur/internal/service/users"
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) setGradesEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/save",
		h.mw(h.handleGradesSave),
	).Methods("POST")

	rtr.HandleFunc(url+"/by_month_and_subject",
		h.mw(h.handleGradesByMonthAndSubject),
	).Methods("POST")

	rtr.HandleFunc(url+"/by_month_and_user",
		h.mw(h.handleGetUserGradesByMonth),
	).Methods("POST")

	rtr.HandleFunc(url+"/by_semester_and_user",
		h.mw(h.handleGetUserGradesBySemester),
	).Methods("POST")

	rtr.HandleFunc(url+"/res_by_subject",
		h.mw(h.handleGetResGradesBySubject),
	).Methods("POST")

	rtr.HandleFunc(url+"/res_by_user",
		h.mw(h.handleGetResUserGrades),
	).Methods("POST")

}

func (h *Handler) handleGradesSave(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, users.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	var in []*models.GradeToSave
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	if err := h.gradesService.Save(r.Context(), in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}

type getGradesByMonthAndSubjectIn struct {
	Month     int8 `json:"month"`
	SubjectId int  `json:"subject_id"`
	Course    int8 `json:"course"`
}
type getGradesByMonthAndSubjectOut struct {
	Days        []int8              `json:"days"`
	Users       []grades.MinUser    `json:"users"`
	Grades      [][]models.MinGrade `json:"grades"`
	SubjectName string              `json:"subject_name"`
}

func (h *Handler) handleGradesByMonthAndSubject(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, users.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	var in getGradesByMonthAndSubjectIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	res, err := h.gradesService.GetByMonthAndSubject(r.Context(), in.Month, in.SubjectId, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	subjectName, err := h.subjectService.GetSubject(r.Context(), in.SubjectId)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&getGradesByMonthAndSubjectOut{
		Days:        res.Days,
		Users:       res.Users,
		Grades:      res.Grades,
		SubjectName: subjectName,
	})
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

type getUserGradesByMonthIn struct {
	Month  int8 `json:"month"`
	Course int8 `json:"course"`
}

func (h *Handler) handleGetUserGradesByMonth(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, users.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	var in getUserGradesByMonthIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	out, err := h.gradesService.GetUserGradesByMonth(r.Context(), login, in.Month, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(out)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

type getUserGradesBySemesterIn struct {
	Semester int8 `json:"semester"`
	Course   int8 `json:"course"`
}

func (h *Handler) handleGetUserGradesBySemester(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, users.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	var in getUserGradesBySemesterIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	out, err := h.gradesService.GetUserGradesBySemester(r.Context(), login, in.Semester, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(out)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

type getResGradesBySubject struct {
	SubjectId int  `json:"subject_id"`
	Course    int8 `json:"course"`
}

func (h *Handler) handleGetResGradesBySubject(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, users.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	var in getResGradesBySubject
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	out, err := h.gradesService.GetResultGradesBySubject(r.Context(), in.SubjectId, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(out)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

type getResUserGradesIn struct {
	Course int8 `json:"course"`
}

func (h *Handler) handleGetResUserGrades(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, users.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	var in getResUserGradesIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	out, err := h.gradesService.GetResultGradesByUser(r.Context(), login, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(out)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}
