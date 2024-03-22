package httpServer

import (
	"eljur/internal/domain/models"
	subjectsService "eljur/internal/service/subjects"
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) setSubjectsEndpoints(rtr *mux.Router, url string) {
	studentUrl = url

	rtr.HandleFunc(url+"/get_by_semester",
		h.mw(h.handleSubjectsGetBySemester),
	).Methods("POST")

	rtr.HandleFunc(url+"/save",
		h.mw(h.handleSubjectsSave),
	).Methods("POST")
}

type getSubjectsBySemesterIn struct {
	Semester int8 `json:"semester"`
	Course   int8 `json:"course"`
}

func (h *Handler) handleSubjectsGetBySemester(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	var in getSubjectsBySemesterIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	subjects, err := h.subjectService.GetBySemester(in.Semester, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
	}
	resp, err := json.Marshal(subjects)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

func (h *Handler) handleSubjectsSave(w http.ResponseWriter, r *http.Request) {
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

	var in []subjectsService.SaveSubjectIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	if err := h.subjectService.Save(in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
	}
}
