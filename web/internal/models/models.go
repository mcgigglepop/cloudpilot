package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          	int
	FirstName    	string
	LastName    	string
	Email    			string
	Password 			string
	Confirmed 		bool
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
	OrganizationId int
}

// Organization is the organization model
type Organization struct {
	ID          	int
	Name    	    string
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}
