package httpServer

import (
	"eljur/internal/domain/models"
	"eljur/internal/service/grades"
	"eljur/internal/service/subjects"
	"eljur/internal/service/users"
	"eljur/pkg/AuthClient"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log/slog"
	"net/http"
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

	rtr.HandleFunc("/admin", h.handleAdmin)

	return rtr
}

func (h *Handler) handleAdmin(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		h.l.Warn(err.Error())
		return
	}
	login, err := h.auth.ParseToken(r.Context(), token)
	if err != nil {
		h.l.Warn(err.Error())
		return
	}
	perm, err := h.auth.GetPermission(r.Context(), login)
	if err != nil {
		h.l.Warn(err.Error())
		return
	}
	if perm < models.PermAdmin {
		h.l.Warn(fmt.Sprintf("unauthorized user"))
		return
	}
	t, err := template.ParseFiles("web/templates/admin.html")
	if err != nil {
		h.l.Warn(err.Error())
	}
	if err := t.Execute(w, nil); err != nil {
		h.l.Warn(err.Error())
	}
}

func redirectToLogin(w http.ResponseWriter, status int) {

}

func getToken(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
