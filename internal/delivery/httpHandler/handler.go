package httpHandler

import (
	"context"
	"eljur/internal/pkg/metric"
	"eljur/internal/service/grades"
	schedules "eljur/internal/service/schedule"
	"eljur/internal/service/subjects"
	"eljur/internal/service/users"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type AuthService interface {
	ParseToken(ctx context.Context, token string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)
	GetPermission(ctx context.Context, login string) (int32, error)
}

type Handler struct {
	l               *slog.Logger
	gradesService   *grades.GradeService
	usersService    *users.UserService
	subjectService  *subjects.SubjectService
	scheduleService *schedules.ScheduleService
	auth            AuthService
}

func NewHandler(l *slog.Logger,
	gradesService *grades.GradeService,
	usersService *users.UserService,
	subjectService *subjects.SubjectService,
	scheduleService *schedules.ScheduleService,
	auth AuthService,
	timeOut time.Duration) *Handler {
	TimeOut = timeOut
	return &Handler{
		l:               l,
		gradesService:   gradesService,
		usersService:    usersService,
		subjectService:  subjectService,
		scheduleService: scheduleService,
		auth:            auth,
	}
}

var TimeOut = time.Second * 10

func (h *Handler) GetMuxRouter() *mux.Router {
	rtr := mux.NewRouter()
	h.setEndpoints(rtr)
	return rtr
}

func (h *Handler) setEndpoints(rtr *mux.Router) {
	rtr.HandleFunc("/", h.handleIndex).Methods("GET")
	rtr.HandleFunc("/login", h.mw(h.loginUser))

	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))

	h.setAdminEndpoints(rtr, "/admin")
	h.setStudentEndpoints(rtr, "/student")
	h.setGradesEndpoints(rtr, "/grades")
	h.setUsersEndpoints(rtr, "/users")
	h.setSubjectsEndpoints(rtr, "/subjects")
	h.setScheduleEndpoints(rtr, "/schedule")
}

type Message struct {
	Mess string
}

type Login struct {
	Login string
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	redirect(w, "/student/grades")
}

func (h *Handler) httpErr(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	h.l.Warn(err.Error())
}

func (h *Handler) renderTemplate(w http.ResponseWriter, data any, filenames ...string) {
	metric.HandleRender()

	for i := range filenames {
		filenames[i] = fmt.Sprintf("web/templates/%s", filenames[i])
	}

	tmp, err := template.ParseFiles(filenames...)
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
