package services

import (
	"context"
	"testing"
	"time"

	"db-gateway-service/proto"
	users "db-gateway-service/sql/user-service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	
	db := sqlx.NewDb(mockDB, "postgres")
	return db, mock
}

func TestUserService_CreateUser(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	req := &proto.CreateUserRequest{
		FullName:      "John Doe",
		Email:         "john@example.com",
		Password:      "password123",
		PhoneNumber:   "1234567890",
		Sex:           "M",
		City:          "New York",
		StateProvince: "NY",
		PostalCode:    "10001",
		CountryCode:   "US",
		Locale:        "en-US",
		Timezone:      "America/New_York",
		UtcOffset:     -5,
	}

	now := time.Now()
	
	// Setup mock expectations
	mock.ExpectQuery(`INSERT INTO USERS`).
		WithArgs(
			req.FullName, req.Email, req.Password,
			req.PhoneNumber, req.Sex, req.City, req.StateProvince,
			req.PostalCode, req.CountryCode, req.Locale, req.Timezone, int(req.UtcOffset),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	// Execute
	resp, err := service.CreateUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int32(1), resp.User.Id)
	assert.Equal(t, req.FullName, resp.User.FullName)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUserByID(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	userID := int32(1)
	now := time.Now()
	phoneNumber := "1234567890"
	sex := "M"
	city := "New York"
	stateProvince := "NY"
	postalCode := "10001"
	countryCode := "US"
	locale := "en-US"
	timezone := "America/New_York"
	utcOffset := -5

	// Setup mock expectations
	mock.ExpectQuery(`SELECT .+ FROM USERS WHERE id = \$1`).
		WithArgs(int(userID)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "full_name", "email", "phone_number", "sex", "city",
			"state_province", "postal_code", "country_code", "locale",
			"timezone", "utc_offset", "created_at", "updated_at",
		}).AddRow(
			userID, "John Doe", "john@example.com", phoneNumber, sex, city,
			stateProvince, postalCode, countryCode, locale,
			timezone, utcOffset, now, now,
		))

	// Execute
	req := &proto.GetUserByIDRequest{Id: userID}
	resp, err := service.GetUserByID(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, userID, resp.User.Id)
	assert.Equal(t, "John Doe", resp.User.FullName)
	assert.Equal(t, "john@example.com", resp.User.Email)
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetAllUsers(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	now := time.Now()
	phoneNumber := "1234567890"
	sex := "M"
	city := "New York"
	stateProvince := "NY"
	postalCode := "10001"
	countryCode := "US"
	locale := "en-US"
	timezone := "America/New_York"
	utcOffset := -5

	// Setup mock expectations
	mock.ExpectQuery(`SELECT .+ FROM USERS ORDER BY created_at DESC`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "full_name", "email", "phone_number", "sex", "city",
			"state_province", "postal_code", "country_code", "locale",
			"timezone", "utc_offset", "created_at", "updated_at",
		}).
			AddRow(1, "John Doe", "john@example.com", phoneNumber, sex, city,
				stateProvince, postalCode, countryCode, locale,
				timezone, utcOffset, now, now).
			AddRow(2, "Jane Smith", "jane@example.com", phoneNumber, sex, city,
				stateProvince, postalCode, countryCode, locale,
				timezone, utcOffset, now, now))

	// Execute
	req := &proto.GetAllUsersRequest{}
	resp, err := service.GetAllUsers(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, int32(1), resp.Users[0].Id)
	assert.Equal(t, "John Doe", resp.Users[0].FullName)
	assert.Equal(t, int32(2), resp.Users[1].Id)
	assert.Equal(t, "Jane Smith", resp.Users[1].FullName)
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_UpdateUser(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	req := &proto.UpdateUserRequest{
		Id:            1,
		FullName:      "John Updated",
		Password:      "newpassword",
		PhoneNumber:   "9876543210",
		Sex:           "M",
		City:          "Los Angeles",
		StateProvince: "CA",
		PostalCode:    "90001",
		CountryCode:   "US",
		Locale:        "en-US",
		Timezone:      "America/Los_Angeles",
		UtcOffset:     -8,
	}

	now := time.Now()

	// Setup mock expectations
	mock.ExpectQuery(`UPDATE USERS .+ WHERE id = \$12 RETURNING`).
		WithArgs(
			req.FullName, req.Password, req.PhoneNumber, req.Sex,
			req.City, req.StateProvince, req.PostalCode, req.CountryCode,
			req.Locale, req.Timezone, int(req.UtcOffset), int(req.Id),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "full_name", "email", "phone_number", "sex", "city",
			"state_province", "postal_code", "country_code", "locale",
			"timezone", "utc_offset", "created_at", "updated_at",
		}).AddRow(
			req.Id, req.FullName, "john@example.com", req.PhoneNumber, req.Sex, req.City,
			req.StateProvince, req.PostalCode, req.CountryCode, req.Locale,
			req.Timezone, req.UtcOffset, now, now,
		))

	// Execute
	resp, err := service.UpdateUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, req.Id, resp.User.Id)
	assert.Equal(t, req.FullName, resp.User.FullName)
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_DeleteUser(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	userID := int32(1)

	// Setup mock expectations
	mock.ExpectExec(`DELETE FROM USERS WHERE id = \$1`).
		WithArgs(int(userID)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	req := &proto.DeleteUserRequest{Id: userID}
	resp, err := service.DeleteUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Message, "deleted successfully")
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_VerifyUser(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	email := "john@example.com"
	password := "password123"
	now := time.Now()

	// Setup mock expectations for successful verification
	mock.ExpectQuery(`SELECT .+ FROM USERS WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "full_name", "email", "password", "phone_number", "sex", "city",
			"state_province", "postal_code", "country_code", "locale",
			"timezone", "utc_offset", "created_at", "updated_at",
		}).AddRow(
			1, "John Doe", email, password, nil, nil, nil,
			nil, nil, nil, nil,
			nil, nil, now, now,
		))

	// Execute
	req := &proto.VerifyUserRequest{
		Email:    email,
		Password: password,
	}
	resp, err := service.VerifyUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.NotNil(t, resp.User)
	assert.Equal(t, email, resp.User.Email)
	assert.Empty(t, resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_VerifyUser_InvalidPassword(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := users.NewRepository(db)
	service := NewUserService(repo)

	// Test data
	email := "john@example.com"
	password := "wrongpassword"
	correctPassword := "password123"
	now := time.Now()

	// Setup mock expectations
	mock.ExpectQuery(`SELECT .+ FROM USERS WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "full_name", "email", "password", "phone_number", "sex", "city",
			"state_province", "postal_code", "country_code", "locale",
			"timezone", "utc_offset", "created_at", "updated_at",
		}).AddRow(
			1, "John Doe", email, correctPassword, nil, nil, nil,
			nil, nil, nil, nil,
			nil, nil, now, now,
		))

	// Execute
	req := &proto.VerifyUserRequest{
		Email:    email,
		Password: password,
	}
	resp, err := service.VerifyUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
	assert.Nil(t, resp.User)
	assert.Equal(t, "Invalid credentials", resp.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
