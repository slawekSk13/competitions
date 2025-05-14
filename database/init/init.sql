-- Drop existing tables if they exist
DROP TABLE IF EXISTS competition_participants;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS competitions;

-- Create tables
CREATE TABLE competitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE participants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE competition_participants (
    competition_id INTEGER NOT NULL,
    participant_id INTEGER NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (competition_id, participant_id),
    FOREIGN KEY (competition_id) REFERENCES competitions(id) ON DELETE CASCADE,
    FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
);

-- Insert sample data
INSERT INTO competitions (name, description, date, location) VALUES
('Summer Athletics Championship', 'Annual athletics event featuring track and field competitions.', '2025-07-15', 'Central Stadium'),
('Winter Swimming Tournament', 'Indoor swimming competition for all age categories.', '2025-12-10', 'Aquatic Center'),
('Chess Masters Championship', 'International chess tournament for professional players.', '2025-09-05', 'Grand Hotel Conference Hall');

-- Insert participants
INSERT INTO participants (name, email) VALUES
('John Smith', 'john.smith@example.com'),
('Anna Johnson', 'anna.johnson@example.com'),
('Michael Brown', 'michael.brown@example.com'),
('Emily Davis', 'emily.davis@example.com'),
('Robert Wilson', 'robert.wilson@example.com'),
('Sarah Thompson', 'sarah.thompson@example.com');

-- Insert competition participants
INSERT INTO competition_participants (competition_id, participant_id, registration_date) VALUES
(1, 1, '2025-05-10'),
(1, 2, '2025-05-11'),
(2, 3, '2025-11-05'),
(2, 4, '2025-11-06'),
(3, 5, '2025-08-15'),
(3, 6, '2025-08-16'),
-- Add some participants to multiple competitions
(1, 3, '2025-05-12'),
(2, 1, '2025-11-01'),
(3, 2, '2025-08-10');
