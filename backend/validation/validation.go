package validation

import (
	"errors"
	"regexp"
	"time"
)

type Competition struct {
	ID          int       
	Name        string    
	Description string    
	Date        time.Time 
	Location    string    
}

type Participant struct {
	ID    int       
	Name  string    
	Email string    
}

// ParseDate parses a date string in format "YYYY-MM-DD"
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ValidateCompetition validates competition data
func ValidateCompetition(c *Competition) error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if len(c.Name) > 255 {
		return errors.New("name is too long (maximum 255 characters)")
	}

	if c.Date.IsZero() {
		return errors.New("date is required")
	}

	if c.Location == "" {
		return errors.New("location is required")
	}

	if len(c.Location) > 255 {
		return errors.New("location is too long (maximum 255 characters)")
	}

	return nil
}

// ValidateParticipant validates participant data
func ValidateParticipant(p *Participant) error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	if len(p.Name) > 255 {
		return errors.New("name is too long (maximum 255 characters)")
	}

	if p.Email == "" {
		return errors.New("email is required")
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}

	return nil
}
