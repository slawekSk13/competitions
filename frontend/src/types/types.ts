export interface Competition {
  id: number;
  name: string;
  description: string;
  date: string;
  location: string;
  created_at?: string;
  updated_at?: string;
}

export interface Participant {
  id: number;
  name: string;
  email: string;
  created_at?: string;
  updated_at?: string;
}

export interface ParticipantFormData {
  name: string;
  email: string;
}

export interface CompetitionParticipant {
  competition_id: number;
  participant_id: number;
  registration_date: string;
  created_at?: string;
  updated_at?: string;
}

export interface CompetitionFormData extends Omit<Competition, "id"> {}
