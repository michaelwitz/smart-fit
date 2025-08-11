package checkins

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// CheckIn represents a check-in record in the database
type CheckIn struct {
	ID           int        `db:"id"`
	UserID       int        `db:"user_id"`
	CheckInType  string     `db:"check_in_type"`
	Notes        *string    `db:"notes"`
	Location     *string    `db:"location"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// Repository handles check-in database operations
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new check-in repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateCheckIn creates a new check-in record
func (r *Repository) CreateCheckIn(checkIn *CheckIn) error {
	query := `
		INSERT INTO CHECK_INS (user_id, check_in_type, notes, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query, checkIn.UserID, checkIn.CheckInType, checkIn.Notes, checkIn.Location,
	).Scan(&checkIn.ID, &checkIn.CreatedAt, &checkIn.UpdatedAt)
}

// GetCheckInByID retrieves a check-in by ID
func (r *Repository) GetCheckInByID(id int) (*CheckIn, error) {
	var checkIn CheckIn
	query := `
		SELECT id, user_id, check_in_type, notes, location, created_at, updated_at
		FROM CHECK_INS 
		WHERE id = $1`

	err := r.db.Get(&checkIn, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("check-in not found")
		}
		return nil, err
	}

	return &checkIn, nil
}

// GetCheckInsByUser retrieves check-ins for a specific user
func (r *Repository) GetCheckInsByUser(userID int, limit, offset int) ([]CheckIn, error) {
	var checkIns []CheckIn
	query := `
		SELECT id, user_id, check_in_type, notes, location, created_at, updated_at
		FROM CHECK_INS 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	err := r.db.Select(&checkIns, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return checkIns, nil
}

// UpdateCheckIn updates an existing check-in
func (r *Repository) UpdateCheckIn(checkIn *CheckIn) error {
	query := `
		UPDATE CHECK_INS 
		SET check_in_type = $1, notes = $2, location = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, user_id, check_in_type, notes, location, created_at, updated_at`

	return r.db.QueryRow(
		query, checkIn.CheckInType, checkIn.Notes, checkIn.Location, checkIn.ID,
	).Scan(
		&checkIn.ID, &checkIn.UserID, &checkIn.CheckInType, &checkIn.Notes, &checkIn.Location,
		&checkIn.CreatedAt, &checkIn.UpdatedAt,
	)
}

// DeleteCheckIn deletes a check-in by ID
func (r *Repository) DeleteCheckIn(id int) error {
	query := `DELETE FROM CHECK_INS WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("check-in not found")
	}

	return nil
}

// GetCheckInsByType retrieves check-ins by type for a user
func (r *Repository) GetCheckInsByType(userID int, checkInType string, limit int) ([]CheckIn, error) {
	var checkIns []CheckIn
	query := `
		SELECT id, user_id, check_in_type, notes, location, created_at, updated_at
		FROM CHECK_INS 
		WHERE user_id = $1 AND check_in_type = $2
		ORDER BY created_at DESC
		LIMIT $3`

	err := r.db.Select(&checkIns, query, userID, checkInType, limit)
	if err != nil {
		return nil, err
	}

	return checkIns, nil
}

// GetCheckInsByDateRange retrieves check-ins within a date range
func (r *Repository) GetCheckInsByDateRange(userID int, startDate, endDate time.Time) ([]CheckIn, error) {
	var checkIns []CheckIn
	query := `
		SELECT id, user_id, check_in_type, notes, location, created_at, updated_at
		FROM CHECK_INS 
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC`

	err := r.db.Select(&checkIns, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return checkIns, nil
}
