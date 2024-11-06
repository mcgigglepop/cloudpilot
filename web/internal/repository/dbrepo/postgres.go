package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"
	"math/rand"

	"github.com/mcgigglepop/cloudpilot/web/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertOrganization inserts a organization into the database
func (m *postgresDBRepo) InsertOrganization(org models.Organization) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into organizations (name, created_at, updated_at) 
			values ($1, $2, $3) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		org.Name,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// CreateUser creates a user in the database
func (m *postgresDBRepo) CreateUser(u models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the email already exists
	var id int
	
	row := m.DB.QueryRowContext(ctx, `select id from users where email = $1`, u.Email)
	_ = row.Scan(&id)

	if id != 0 {
		return "", errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Insert the new user into the database
	stmt := `
		insert into users (first_name, last_name, email, password_hash, confirmed, created_at, updated_at, organization_id) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		hashedPassword,
		false,
		time.Now(),
		time.Now(),
		u.OrganizationId,
	)

	// error handle
	if err != nil {
		log.Println(err)
		return "", err
	}

	// return username
	return u.Email, nil
}

// AddCloudAccount creates a cloud account in the database
func (m *postgresDBRepo) AddCloudAccount(c models.CloudAccount) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into cloud_accounts (account_number, provider_name, account_type, organization_id, account_status, in_use, created_at, updated_at) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		c.AccountNumber,
		c.ProviderName, 
		c.AccountType,
		c.OrganizationID,
		"unavailable",
		false,
		time.Now(),
		time.Now(),
	)

	// error handle
	if err != nil {
		log.Println(err)
		return err
	}
	
	// return username
	return nil
}

// GetUserByID gets a room by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	query := `
		select user_id, organization_id from users where user_id = $1	
`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.OrganizationId,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, `select user_id, password_hash from users where email = $1`, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// RandomOrgName creates a random organization name
func (m *postgresDBRepo) RandomOrgName() string {
	names := []string{"SunshineClouds", "MoonlightDev", "StarrySkies"} // Add more names as desired
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}