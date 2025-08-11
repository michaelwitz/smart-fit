package users

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// User represents a user in the database
type User struct {
	ID            int        `db:"id"`
	FullName      string     `db:"full_name"`
	Email         string     `db:"email"`
	Password      string     `db:"password"`
	PhoneNumber   *string    `db:"phone_number"`
	Sex           *string    `db:"sex"`
	City          *string    `db:"city"`
	StateProvince *string    `db:"state_province"`
	PostalCode    *string    `db:"postal_code"`
	CountryCode   *string    `db:"country_code"`
	Locale        *string    `db:"locale"`
	Timezone      *string    `db:"timezone"`
	UtcOffset     *int       `db:"utc_offset"`
	LastActive    *time.Time `db:"last_active"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

// Repository handles user database operations
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user
func (r *Repository) CreateUser(user *User) error {
	query := `
		INSERT INTO USERS (full_name, email, password, phone_number, sex, 
		                  city, state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query, user.FullName, user.Email, user.Password,
		user.PhoneNumber, user.Sex, user.City, user.StateProvince,
		user.PostalCode, user.CountryCode, user.Locale, user.Timezone, user.UtcOffset,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(id int) (*User, error) {
	var user User
	query := `
		SELECT id, full_name, email, phone_number, sex, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		WHERE id = $1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (r *Repository) GetAllUsers() ([]User, error) {
	var users []User
	query := `
		SELECT id, full_name, email, phone_number, sex, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		ORDER BY created_at DESC`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *Repository) UpdateUser(user *User) error {
	query := `
		UPDATE USERS 
		SET full_name = $1, password = $2, phone_number = $3, sex = $4, 
		    city = $5, state_province = $6, postal_code = $7, country_code = $8, locale = $9, timezone = $10, 
		    utc_offset = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		RETURNING id, full_name, email, phone_number, sex, city, 
		          state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at`

	return r.db.QueryRow(
		query, user.FullName, user.Password, user.PhoneNumber, user.Sex,
		user.City, user.StateProvince, user.PostalCode, user.CountryCode, user.Locale, user.Timezone, 
		user.UtcOffset, user.ID,
	).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PhoneNumber, &user.Sex,
		&user.City, &user.StateProvince, &user.PostalCode,
		&user.CountryCode, &user.Locale, &user.Timezone, &user.UtcOffset, 
		&user.CreatedAt, &user.UpdatedAt,
	)
}

// DeleteUser deletes a user by ID
func (r *Repository) DeleteUser(id int) error {
	query := `DELETE FROM USERS WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// VerifyUser verifies user credentials
func (r *Repository) VerifyUser(email, password string) (*User, error) {
	var user User
	query := `
		SELECT id, full_name, email, password, phone_number, sex, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		WHERE email = $1`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Note: Password verification should be done by the calling service
	// This method just retrieves the user data

	return &user, nil
}

// UpsertUser creates or updates a user
func (r *Repository) UpsertUser(user *User) error {
	// Try to get existing user
	existingUser, err := r.GetUserByID(user.ID)
	if err != nil {
		// User doesn't exist, create new one
		return r.CreateUser(user)
	}

	// User exists, update it
	user.ID = existingUser.ID
	return r.UpdateUser(user)
}

// UpdateUserPartial updates specific user fields
func (r *Repository) UpdateUserPartial(id int, updates map[string]interface{}) (*User, error) {
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	for field, value := range updates {
		if value != nil {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	// Always update the updated_at field
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE USERS SET %s WHERE id = $%d
		RETURNING id, full_name, email, phone_number, sex, city, 
		          state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var user User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PhoneNumber, &user.Sex,
		&user.City, &user.StateProvince, &user.PostalCode,
		&user.CountryCode, &user.Locale, &user.Timezone, &user.UtcOffset, 
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
