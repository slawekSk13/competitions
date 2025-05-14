import axios from "axios";
import {
  Competition,
  Participant,
  CompetitionParticipant,
} from "../types/types";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  timeout: 5000,
  headers: {
    "Content-Type": "application/json",
  },
});

export const CompetitionAPI = {
  getAll: () => api.get<Competition[]>("/competitions"),
  getById: (id: number) => api.get<Competition>(`/competitions/${id}`),
  create: (data: Competition) => api.post<Competition>("/competitions", data),
  update: (id: number, data: Competition) =>
    api.put<Competition>(`/competitions/${id}`, data),
  delete: (id: number) => api.delete(`/competitions/${id}`),
};

export const ParticipantAPI = {
  getAll: (competitionId?: number) =>
    api.get<Participant[]>("/participants", {
      params: { competition_id: competitionId },
    }),
  getById: (id: number) => api.get<Participant>(`/participants/${id}`),
  getCompetitions: (id: number) =>
    api.get<Competition[]>(`/participants/${id}/competitions`),
  create: (data: Participant) => api.post<Participant>("/participants", data),
  addToCompetition: (
    id: number,
    data: { competition_id: number; registration_date: string }
  ) =>
    api.post<CompetitionParticipant>(`/participants/${id}/competitions`, data),
  update: (id: number, data: Participant) =>
    api.put<Participant>(`/participants/${id}`, data),
  removeFromCompetition: (id: number, competitionId: number) =>
    api.delete(`/participants/${id}/competitions/${competitionId}`),
  delete: (id: number) => api.delete(`/participants/${id}`),
};
