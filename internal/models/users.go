package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (*User, error)
	Exists(id int) (bool, error)
}

func (m *UserModel) Insert(name, email, password string) error {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(HashedPassword))

	if err != nil {
		var mySQLError *mysql.MySQLError

		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (*User, error) {
	stmt := `SELECT id, hashed_password, Name FROM users WHERE email = ?`

	u := &User{}

	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.HashedPassword, &u.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, ErrInvalidCredentials
		} else {
			return u, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return u, ErrInvalidCredentials
		} else {
			return u, err
		}
	}

	return u, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
