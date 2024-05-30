package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/internal/service/schedule"
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) setScheduleEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/save",
		h.mw(h.handleScheduleSave),
	).Methods("POST")

	rtr.HandleFunc(url+"/get_all",
		h.mw(h.handleScheduleGetAll),
	).Methods("POST")

	rtr.HandleFunc(url+"/get_actual",
		h.mw(h.handleScheduleGetActual),
	).Methods("POST")
}

func (h *Handler) handleScheduleGetActual(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	res, err := h.scheduleService.GetActualSchedule(ctx, login)
	if err != nil {
		h.l.Warn(tr.Trace(err).Error())
	}
	out, err := json.Marshal(res)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(out)
}

func (h *Handler) handleScheduleGetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	res, err := h.scheduleService.GetAll(ctx)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(res)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(out)
}

func (h *Handler) handleScheduleSave(w http.ResponseWriter, r *http.Request) {
	_, ok := h.authenticate(r, models.PermAdmin)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	var in []schedules.LessonToSave
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	if err := h.scheduleService.Save(ctx, in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}
