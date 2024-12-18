package mocks

import (
	"time"

	"snippetbox.fepg.org/internal/models"
)

type UserModel struct{}

var testUserData = struct {
	Id       int
	Name     string
	Email    string
	Password string
	Created  time.Time
}{
	Id:       1,
	Name:     "Alice",
	Email:    "alice@example.com",
	Password: "pa$$word",
	Created:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == testUserData.Email && password == testUserData.Password {
		return testUserData.Id, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case testUserData.Id:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) GetUserData(id int) (*models.UserData, error) {
	if id == testUserData.Id {
		userData := &models.UserData{
			Name:    testUserData.Name,
			Email:   testUserData.Email,
			Created: testUserData.Created,
		}
		return userData, nil
	}
	return nil, models.ErrNoRecord
}

func (m *UserModel) AuthenticateUsingID(id int, password string) error {
	if id == testUserData.Id && password == testUserData.Password {
		return nil
	}
	return models.ErrInvalidCredentials
}

func (m *UserModel) UpdateName(id int, name, password string) error {
	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	testUserData.Name = name

	return nil
}

func (m *UserModel) UpdateEmail(id int, email, password string) error {
	err := m.AuthenticateUsingID(id, password)
	if err != nil {
		return err
	}

	testUserData.Email = email

	return nil
}

func (m *UserModel) UpdatePassword(id int, oldPassword, newPassword string) error {
	err := m.AuthenticateUsingID(id, oldPassword)
	if err != nil {
		return err
	}

	testUserData.Password = newPassword

	return nil
}
