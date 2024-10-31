package handlers

import (
	"net/http"
	"log"
	"github.com/mcgigglepop/cloudpilot/web/internal/config"
	"github.com/mcgigglepop/cloudpilot/web/internal/driver"
	"github.com/mcgigglepop/cloudpilot/web/internal/models"
	"github.com/mcgigglepop/cloudpilot/web/internal/render"
	"github.com/mcgigglepop/cloudpilot/web/internal/repository"
	"github.com/mcgigglepop/cloudpilot/web/internal/repository/dbrepo"
	"github.com/mcgigglepop/cloudpilot/web/internal/forms"
)

// repository used by the handlers
var Repo *Repository

// the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Index is the home/index page handler
func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index.page.tmpl", &models.TemplateData{})
}

// Login is the login page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
}

// PostLogin handles logging the user in
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	// parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// get email, password
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// required fields
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	// check form validation
	if !form.Valid() {
		log.Println("Invalid form")
		m.App.Session.Put(r.Context(), "error", "invalid form")
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	// authenticate the user
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println("Invalid login credentials")
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// put the user id in the session and log in user
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// Register is the register page handler
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{})
}

// PostRegister is the post - register page handler
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	// parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// user model
	user := models.User{
		FirstName: r.Form.Get("firstname"),
		LastName:  r.Form.Get("lastname"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
		Confirmed: false,
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	// validate the form
	if !form.Valid() {
		log.Println("Invalid form")
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	// create the user, returning the email for session storage
	email, err := m.DB.CreateUser(user)

	// validate the CreateUser call
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// store session data and re-direct to confirm-email
	m.App.Session.Put(r.Context(), "flash", "registered successfully")
	m.App.Session.Put(r.Context(), "email", email)
	http.Redirect(w, r, "/confirm-email", http.StatusSeeOther)
}

// Dashboard is the dashboard page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{})
}