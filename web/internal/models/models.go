package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          	int
	Username    	string
	Password 	string
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}