package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Name    string
	Email   string
	Created time.Time
}

type User struct {
	ID int
	//Name           string
	//Email          string
	HashedPassword []byte
	//Created        time.Time
	*UserData
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	GetUserData(id int) (*UserData, error)
	UpdateName(id int, name, password string) error
	UpdateEmail(id int, email, password string) error
	UpdatePassword(id int, oldPasswd, newPasswd string) error
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	const stmt = `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
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
	var id int
	var hashedPassword []byte

	const stmt = "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UserModel) AuthenticateUsingID(id int, password string) error {
	var hashed_password []byte

	stmt := "SELECT hashed_password FROM users WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	const stmt = "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrInvalidCredentials
		}

		return false, nil
	}

	return exists, nil
}

func (m *UserModel) GetUserData(id int) (*UserData, error) {
	userData := &UserData{}

	const stmt = "SELECT name, email, created FROM users WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&userData.Name, &userData.Email, &userData.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return userData, nil
}

func (m *UserModel) UpdateName(id int, name, password string) error {
	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET name = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, name, id)

	return err
}

func (m *UserModel) UpdateEmail(id int, email, password string) error {
	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET email = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, email, id)

	return err
}

func (m *UserModel) UpdatePassword(id int, oldPasswd, newPasswd string) error {
	err := m.AuthenticateUsingID(id, oldPasswd)
	if err != nil {
		return err
	}

	newHashedPasswd, err := bcrypt.GenerateFromPassword([]byte(newPasswd), 12)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, string(newHashedPasswd), id)

	return err
}
