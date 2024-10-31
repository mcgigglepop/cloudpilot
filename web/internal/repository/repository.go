package repository

import (
	"github.com/mcgigglepop/cloudpilot/web/internal/models"
)

type DatabaseRepo interface {
	CreateUser(u models.User) (string, error)
	Authenticate(email, password string) (int, string, error)
}