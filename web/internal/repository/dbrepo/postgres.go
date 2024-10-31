package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mcgigglepop/cloudpilot/web/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// CreateUser creates a user in the database
func (m *postgresDBRepo) CreateUser(u models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the email already exists
	var id int
	
	row := m.DB.QueryRowContext(ctx, `select id from users where username = $1`, u.Username)
	_ = row.Scan(&id)

	if id != 0 {
		return "", errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Insert the new user into the database
	stmt := `
		insert into users (username, password_hash, confirmed, created_at, updated_at) 
		values ($1, $2, $3, $4, $5)
	`
	_, err = m.DB.ExecContext(ctx, stmt,
		u.Username,
		hashedPassword,
		false,
		time.Now(),
		time.Now(),
	)

	// error handle
	if err != nil {
		log.Println(err)
		return "", err
	}

	// return username
	return u.Username, nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(username, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, `select user_id, password from users where username = $1`, username)
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