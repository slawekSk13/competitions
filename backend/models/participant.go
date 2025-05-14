package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Participant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompetitionParticipant struct {
	CompetitionID    int       `json:"competition_id"`
	ParticipantID    int       `json:"participant_id"`
	RegistrationDate time.Time `json:"registration_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Participant
func (p *Participant) UnmarshalJSON(data []byte) error {
	type Alias Participant
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// GetAllParticipants retrieves all participants
func GetAllParticipants() ([]Participant, error) {
	rows, err := DB.Query(`
		SELECT id, name, email, created_at, updated_at 
		FROM participants
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []Participant
	for rows.Next() {
		var p Participant
		err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// GetParticipantsByCompetition retrieves all participants for a specific competition
func GetParticipantsByCompetition(competitionID int) ([]Participant, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.name, p.email, p.created_at, p.updated_at 
		FROM participants p
		JOIN competition_participants cp ON p.id = cp.participant_id
		WHERE cp.competition_id = $1
		ORDER BY cp.registration_date DESC
	`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []Participant
	for rows.Next() {
		var p Participant
		err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// GetParticipant retrieves a single participant by ID
func GetParticipant(id int) (Participant, error) {
	var p Participant
	err := DB.QueryRow(`
		SELECT id, name, email, created_at, updated_at 
		FROM participants 
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Email, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return p, errors.New("participant not found")
	}

	return p, err
}

// GetParticipantCompetitions retrieves all competitions for a specific participant
func GetParticipantCompetitions(participantID int) ([]Competition, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.name, c.description, c.date, c.location, c.created_at, c.updated_at
		FROM competitions c
		JOIN competition_participants cp ON c.id = cp.competition_id
		WHERE cp.participant_id = $1
		ORDER BY c.date ASC
	`, participantID)
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

// CreateParticipant adds a new participant to the database
func CreateParticipant(p *Participant) error {
	// Check if email is already used
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM participants WHERE email = $1", p.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already registered")
	}

	err = DB.QueryRow(`
		INSERT INTO participants (name, email)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`, p.Name, p.Email).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	return err
}

// AddParticipantToCompetition adds a participant to a competition
func AddParticipantToCompetition(participantID, competitionID int, registrationDate time.Time) error {
	// Check if the competition exists
	if !CompetitionExists(competitionID) {
		return errors.New("competition does not exist")
	}

	// Check if the participant exists
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM participants WHERE id = $1)", participantID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("participant does not exist")
	}

	// Check if the participant is already registered for this competition
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM competition_participants WHERE participant_id = $1 AND competition_id = $2)", 
		participantID, competitionID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("participant already registered for this competition")
	}

	_, err = DB.Exec(`
		INSERT INTO competition_participants (participant_id, competition_id, registration_date)
		VALUES ($1, $2, $3)
	`, participantID, competitionID, registrationDate)

	return err
}

// UpdateParticipant updates an existing participant
func UpdateParticipant(p *Participant) error {
	// Check if email is already used by another participant
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM participants WHERE email = $1 AND id != $2", p.Email, p.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already registered")
	}

	result, err := DB.Exec(`
		UPDATE participants
		SET name = $2, email = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, p.ID, p.Name, p.Email)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("participant not found")
	}

	return nil
}

// RemoveParticipantFromCompetition removes a participant from a competition
func RemoveParticipantFromCompetition(participantID, competitionID int) error {
	result, err := DB.Exec(`
		DELETE FROM competition_participants 
		WHERE participant_id = $1 AND competition_id = $2
	`, participantID, competitionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("participant not found in competition")
	}

	return nil
}

// DeleteParticipant removes a participant from the database
func DeleteParticipant(id int) error {
	result, err := DB.Exec("DELETE FROM participants WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("participant not found")
	}

	return nil
}
