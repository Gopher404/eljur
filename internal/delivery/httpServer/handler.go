package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/internal/service/grades"
	"eljur/internal/service/subjects"
	"eljur/internal/service/users"
	"eljur/pkg/AuthClient"
	"eljur/pkg/tr"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	l              *slog.Logger
	gradesService  *grades.GradeService
	usersService   *users.UserService
	subjectService *subjects.SubjectService
	auth           *AuthClient.Client
}

func NewHandler(l *slog.Logger,
	gradesService *grades.GradeService,
	usersService *users.UserService,
	subjectService *subjects.SubjectService,
	auth *AuthClient.Client) *Handler {
	return &Handler{
		l:              l,
		gradesService:  gradesService,
		usersService:   usersService,
		subjectService: subjectService,
		auth:           auth,
	}
}

func (h *Handler) GetMuxRouter() *mux.Router {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/",
		h.logHandle(h.handleIndex),
	).Methods("GET")

	rtr.HandleFunc("/admin",
		h.logHandle(h.handleAdmin),
	).Methods("GET")

	rtr.HandleFunc("/login_student",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermStudent, "/")
		}),
	).Methods("GET", "POST")

	rtr.HandleFunc("/login_admin",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermAdmin, "/admin")
		}),
	).Methods("GET", "POST")

	rtr.HandleFunc("/user_grades_by_month",
		h.handleGetUserGradesByMonth,
	).Methods("POST")

	rtr.HandleFunc("/grades_by_month_and_subject",
		h.logHandle(h.getGradesByMonthAndSubject),
	).Methods("POST")

	return rtr
}

type Message struct {
	Mess string
}

type Login struct {
	Login string
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {

	login, ok := h.authenticate(r, models.PermStudent)
	if !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login_student")
		return
	}

	h.renderTemplate(w, "index.html", Login{Login: login})

}

func (h *Handler) handleAdmin(w http.ResponseWriter, r *http.Request) {
	if login, ok := h.authenticate(r, models.PermAdmin); !ok {
		h.l.Info(fmt.Sprintf("unauthorized user %s", login))
		redirect(w, "/login_admin")
		return
	}

	h.renderTemplate(w, "admin.html", nil)
}

type getUserGradesByMonthIn struct {
	Month  int8 `json:"month"`
	Course int8 `json:"course"`
}

func (h *Handler) handleGetUserGradesByMonth(w http.ResponseWriter, r *http.Request) {
	login, ok := h.authenticate(r, models.PermStudent)
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

	out, err := h.gradesService.GetUserGradesByMonth(login, in.Month, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	h.l.Error(fmt.Sprintf("%+v", out))
	resp, err := json.Marshal(out)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
	h.l.Info(string(resp))
	_, _ = w.Write(resp)
}

type getGradesByMonthAndSubjectIn struct {
	Month     int8 `json:"month"`
	SubjectId int  `json:"subject_id"`
	Course    int8 `json:"course"`
}

func (h *Handler) getGradesByMonthAndSubject(w http.ResponseWriter, r *http.Request) {
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

	var in getGradesByMonthAndSubjectIn
	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	out, err := h.gradesService.GetByMonthAndSubject(in.Month, in.SubjectId, in.Course)
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

// tools

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request, perm int32, outUrl string) {
	if r.Method == "GET" {
		h.renderTemplate(w, "login.html", Message{})
	} else {
		if err := r.ParseForm(); err != nil {
			h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
			return
		}

		login := r.Form.Get("login")
		pass := r.Form.Get("pass")

		token, err := h.auth.Login(r.Context(), login, pass)
		if err != nil {
			h.renderTemplate(w, "login.html", Message{Mess: "Неверный логин или пароль"})
			h.l.Warn(tr.Trace(err).Error())
		}

		realPerm, err := h.auth.GetPermission(r.Context(), login)
		if err != nil || realPerm < perm {
			h.renderTemplate(w, "login.html", Message{Mess: "Недостаточно прав"})
			h.l.Warn(tr.Trace(err).Error())
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(time.Hour * 24 * 30),
		})

		redirect(w, outUrl)
	}
}

func (h *Handler) updateGrades(w http.ResponseWriter, r *http.Request) {
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

	var in []models.MinGrade

	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	if err := h.gradesService.UpdateGrades(in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) saveGrades(w http.ResponseWriter, r *http.Request) {
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

	var in []*models.Grade

	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	if err := h.gradesService.SaveGrades(in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) httpErr(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	h.l.Warn(err.Error())
}

func (h *Handler) renderTemplate(w http.ResponseWriter, fileName string, data any) {
	tmp, err := template.ParseFiles(fmt.Sprintf("web/templates/%s", fileName))
	if err != nil {
		h.httpErr(w, err, http.StatusInternalServerError)
	}
	if err := tmp.Execute(w, data); err != nil {
		h.httpErr(w, err, http.StatusInternalServerError)
	}
}

func getToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", tr.Trace(err)
	}
	return c.Value, nil
}

func redirect(w http.ResponseWriter, url string) {
	_, _ = w.Write([]byte(fmt.Sprintf("<script>window.location.replace(%q)</script>", url)))
}
