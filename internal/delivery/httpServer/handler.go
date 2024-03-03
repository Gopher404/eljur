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
	h.setEndpoints(rtr)
	return rtr
}

func (h *Handler) setEndpoints(rtr *mux.Router) {
	rtr.HandleFunc("/",
		h.logHandle(h.handleIndex),
	).Methods("GET")

	h.setAdminEndpoints(rtr, "/admin")

	h.setStudentEndpoints(rtr, "/student")

	h.setGradesEndpoints(rtr, "/grades")

	rtr.HandleFunc("/user_grades_by_month",
		h.handleGetUserGradesByMonth,
	).Methods("POST")

	rtr.HandleFunc("/grades_by_month_and_subject",
		h.logHandle(h.getGradesByMonthAndSubject),
	).Methods("POST")

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
type getGradesByMonthAndSubjectOut struct {
	Days        []int8              `json:"days"`
	Users       []grades.MinUser    `json:"users"`
	Grades      [][]models.MinGrade `json:"grades"`
	SubjectName string              `json:"subject_name"`
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

	res, err := h.gradesService.GetByMonthAndSubject(in.Month, in.SubjectId, in.Course)
	if err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusInternalServerError)
		return
	}

	subjectName, err := h.subjectService.GetSubject(in.SubjectId)
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

const TTL = time.Hour * 24 * 30

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
			return
		}

		realPerm, err := h.auth.GetPermission(r.Context(), login)
		if err != nil || realPerm < perm {
			h.renderTemplate(w, "login.html", Message{Mess: "Недостаточно прав"})
			h.l.Warn(tr.Trace(err).Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(TTL),
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

func (h *Handler) deleteGrades(w http.ResponseWriter, r *http.Request) {
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

	var in []int

	if err := json.Unmarshal(body, &in); err != nil {
		h.httpErr(w, tr.Trace(err), http.StatusBadRequest)
		return
	}

	if err := h.gradesService.DeleteGrades(in); err != nil {
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
