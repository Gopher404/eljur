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
			h.loginUser(w, r, models.PermStudent)
		}),
	).Methods("GET", "POST")

	rtr.HandleFunc("/login_admin",
		h.logHandle(func(w http.ResponseWriter, r *http.Request) {
			h.loginUser(w, r, models.PermAdmin)
		}),
	).Methods("GET", "POST")

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
		redirect(w, "/login_admin")
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

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request, perm int32) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/login.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.l.Warn(err.Error())
		}
		if err := t.Execute(w, struct{ Mess string }{Mess: ""}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.l.Warn(err.Error())
		}
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.l.Warn(err.Error())
			return
		}

		login := r.Form.Get("login")
		pass := r.Form.Get("pass")

		token, err := h.auth.Login(r.Context(), login, pass)
		if err != nil {
			t, err := template.ParseFiles("web/templates/login.html")
			if err != nil {
				h.httpErr(w, err, http.StatusInternalServerError)
			}
			if err := t.Execute(w, struct{ Mess string }{Mess: "Неверный логин или пароль"}); err != nil {
				h.httpErr(w, err, http.StatusInternalServerError)
			}
			h.l.Warn(err.Error())
		}

		realPerm, err := h.auth.GetPermission(r.Context(), login)
		if err != nil || realPerm < perm {
			h.renderTemplate(w, "login.html", Message{Mess: "Недостаточно прав"})
			h.l.Warn(err.Error())
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(time.Hour * 24 * 30),
		})

		redirect(w, "/admin")
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
		return "", err
	}
	return c.Value, nil
}

func redirect(w http.ResponseWriter, url string) {
	_, _ = w.Write([]byte(fmt.Sprintf("<script>window.location.replace(%q)</script>", url)))
}
