package models

import (
	"testing"

	"snippetbox.fepg.org/internal/assert"
)

const (
	validID       = 1
	validEmail    = "alice@example.com"
	validPassword = "pa$$word"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	t.Run("Successful insert", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		err := m.Insert("Alberto Magno", "alberto@exmaple.com", "qqqqwwww")
		assert.NilError(t, err)
	})

	t.Run("Duplicate user", func(t *testing.T) {
		db := newTestDB(t)
		m := UserModel{db}

		err := m.Insert("Alberto Magno", "alice@example.com", "12341234")

		//assert.Equal(t, err, ErrDuplicateEmail)
		assert.ErrorEqual(t, err, ErrDuplicateEmail)
	})

}

func TestUserAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		password  string
		wantId    int
		wantError error
	}{
		{
			name:      "Valid authenticate",
			email:     validEmail,
			password:  validPassword,
			wantId:    1,
			wantError: nil,
		},
		{
			name:      "Empty password",
			email:     validEmail,
			password:  "",
			wantId:    0,
			wantError: ErrInvalidCredentials,
		},
		{
			name:      "Empty email",
			email:     "",
			password:  validPassword,
			wantId:    0,
			wantError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}

			id, err := m.Authenticate(tt.email, tt.password)
			//t.Logf("got %d, wanted %d", id, tt.wantId)

			assert.Equal(t, id, tt.wantId)
			if tt.wantError == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorEqual(t, err, tt.wantError)
			}
		})

	}

}

func TestUserAuthenticateUsingID(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		id        int
		password  string
		wantError error
	}{
		{
			name:      "Valid authentication",
			id:        validID,
			password:  validPassword,
			wantError: nil,
		},
		{
			name:      "Invalid ID",
			id:        0,
			password:  validPassword,
			wantError: ErrInvalidCredentials,
		},
		{
			name:      "Invalid password",
			id:        validID,
			password:  "asdfasdf",
			wantError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			err := m.AuthenticateUsingID(tt.id, tt.password)
			if tt.wantError == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorEqual(t, err, tt.wantError)
			}
		})
	}
}

func TestUserModelGetUserData(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration tests")
	}

	tests := []struct {
		name      string
		id        int
		wantError error
	}{
		{
			name:      "Valid ID",
			id:        validID,
			wantError: nil,
		},
		{
			name:      "Zero ID",
			id:        0,
			wantError: ErrInvalidCredentials,
		},
		{
			name:      "Negative ID",
			id:        -1,
			wantError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}

			//userData, err := m.GetUserData(tt.id)
			_, err := m.GetUserData(tt.id)
			if tt.wantError == nil {
				assert.NilError(t, err)
			}

		})
	}
}

func TestUserModelUpdateName(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration tests")
	}

	tests := []struct {
		name        string
		id          int
		newUserName string
		password    string
		wantError   error
	}{
		{
			name:        "Valid name update",
			id:          validID,
			newUserName: "Susana",
			password:    validPassword,
			wantError:   nil,
		},
		{
			name:        "Zero ID",
			id:          0,
			newUserName: "Susana",
			password:    validPassword,
			wantError:   ErrInvalidCredentials,
		},
		{
			name:        "Invalid password",
			id:          validID,
			newUserName: "Susana",
			password:    "qqqqwwww",
			wantError:   ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}

			err := m.UpdateName(tt.id, tt.newUserName, tt.password)
			if tt.wantError == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorEqual(t, err, tt.wantError)
			}
		})
	}
}

func TestUserModelUpdateEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration tests")
	}

	tests := []struct {
		name      string
		id        int
		newEmail  string
		password  string
		wantError error
	}{
		{
			name:      "Valid email update",
			id:        validID,
			newEmail:  "new_email@example.com",
			password:  validPassword,
			wantError: nil,
		},
		{
			name:      "Invalid ID",
			id:        -1,
			newEmail:  "new_email@example.com",
			password:  validPassword,
			wantError: ErrInvalidCredentials,
		},
		{
			name:      "Invalid password",
			id:        validID,
			newEmail:  "new_email@example.com",
			password:  "11223344",
			wantError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			err := m.UpdateEmail(tt.id, tt.newEmail, tt.password)
			if tt.wantError == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorEqual(t, err, tt.wantError)
			}
		})
	}
}

func TestUserModelUpdatePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration testing")
	}

	const newPassword = "the_new_password"
	const invalidPassword = "invalid_password"

	tests := []struct {
		name        string
		id          int
		oldPassword string
		newPassword string
		wantError   error
	}{
		{
			name:        "Valid password update",
			id:          validID,
			oldPassword: validPassword,
			newPassword: newPassword,
			wantError:   nil,
		},
		{
			name:        "Invalid ID",
			id:          -1,
			oldPassword: validPassword,
			newPassword: newPassword,
			wantError:   ErrInvalidCredentials,
		},
		{
			name:        "Invalid password",
			id:          validID,
			oldPassword: invalidPassword,
			newPassword: newPassword,
			wantError:   ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			err := m.UpdatePassword(tt.id, tt.oldPassword, tt.newPassword)
			if tt.wantError == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorEqual(t, err, tt.wantError)
			}
		})
	}
}
