package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Competition struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Competition
func (c *Competition) UnmarshalJSON(data []byte) error {
	type Alias Competition
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return err
		}
		c.Date = parsedDate
	}
	return nil
}

// GetAllCompetitions retrieves all competitions from the database
func GetAllCompetitions() ([]Competition, error) {
	rows, err := DB.Query(`
		SELECT id, name, description, date, location, created_at, updated_at 
		FROM competitions
		ORDER BY date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitions []Competition
	for rows.Next() {
		var c Competition
		err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Date, &c.Location, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		competitions = append(competitions, c)
	}

	return competitions, nil
}

// GetCompetition retrieves a single competition by ID
func GetCompetition(id int) (Competition, error) {
	var c Competition
	err := DB.QueryRow(`
		SELECT id, name, description, date, location, created_at, updated_at 
		FROM competitions 
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Description, &c.Date, &c.Location, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return c, errors.New("competition not found")
	}

	return c, err
}

// CreateCompetition adds a new competition to the database
func CreateCompetition(c *Competition) error {
	err := DB.QueryRow(`
		INSERT INTO competitions (name, description, date, location)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, c.Name, c.Description, c.Date, c.Location).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	return err
}

// UpdateCompetition updates an existing competition
func UpdateCompetition(c *Competition) error {
	result, err := DB.Exec(`
		UPDATE competitions
		SET name = $2, description = $3, date = $4, location = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, c.ID, c.Name, c.Description, c.Date, c.Location)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("competition not found")
	}

	return nil
}

// DeleteCompetition removes a competition from the database
func DeleteCompetition(id int) error {
	result, err := DB.Exec("DELETE FROM competitions WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("competition not found")
	}

	return nil
}

// CompetitionExists checks if a competition with the given ID exists
func CompetitionExists(id int) bool {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM competitions WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
