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
	h.setEndpoints(rtr)
	return rtr
}

func (h *Handler) setEndpoints(rtr *mux.Router) {
	rtr.HandleFunc("/",
		h.logHandle(h.handleIndex),
	).Methods("GET")

	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))

	h.setAdminEndpoints(rtr, "/admin")
	h.setStudentEndpoints(rtr, "/student")
	h.setGradesEndpoints(rtr, "/grades")
	h.setUsersEndpoints(rtr, "/users")
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

func redirect(w http.ResponseWriter, url string) {
	_, _ = w.Write([]byte(fmt.Sprintf("<script>window.location.replace(%q)</script>", url)))
}
