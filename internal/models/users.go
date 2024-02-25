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

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id FROM users WHERE email = ? AND hashed_password = ?`

	u := &User{}

	err := m.DB.QueryRow(stmt, email, password).Scan(&u.ID)

	if err != nil {
		return 0, ErrInvalidCredentials
	}

	if u.ID < 1 {
		return 0, ErrNoRecord
	}

	return u.ID, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM users WHERE id = ?`

	u := &User{}

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID)

	if err != nil {
		return false, err
	}

	if u.ID < 1 {
		return false, ErrNoRecord
	}

	return true, nil
}
