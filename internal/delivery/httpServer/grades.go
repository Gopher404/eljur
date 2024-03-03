package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) setGradesEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/save",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.gradesSave(w, r)
		}),
	).Methods("GET", "POST")
}

func (h *Handler) gradesSave(w http.ResponseWriter, r *http.Request) {
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

	var in []*models.GradeToSave
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	if err := h.gradesService.Save(in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}
