package httpHandler

import (
	"eljur/pkg/tr"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) setFilesEndpoints(rtr *mux.Router, url string) {
	rtr.HandleFunc(url+"/get",
		h.mw(h.handleFilesGet),
	).Methods("POST")

	rtr.HandleFunc(url+"/save",
		h.mw(h.handleFilesSave),
	).Methods("POST")

	rtr.HandleFunc(url+"/list_dir",
		h.mw(h.handleFilesListDir),
	).Methods("POST")

	rtr.HandleFunc(url+"/create_dir",
		h.mw(h.handleFilesCreateDir),
	).Methods("POST")

}

func (h *Handler) handleFilesGet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	path := r.Form.Get("path")

	file, err := h.fileService.Get(path)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, file.Name, file.ModTime, file.Data)
}

func (h *Handler) handleFilesSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	path := r.Form.Get("path")

	if err := h.fileService.SaveFile(path, r.Body); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleFilesListDir(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	path := r.Form.Get("path")

	files, err := h.fileService.ListDir(path)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(files)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(res)
}

func (h *Handler) handleFilesCreateDir(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}
	path := r.Form.Get("path")

	if err := h.fileService.CreateDir(path); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}
