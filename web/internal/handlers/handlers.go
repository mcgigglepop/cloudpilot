package handlers

import (
	"net/http"
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
