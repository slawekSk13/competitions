import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Competition, Participant } from "../types/types";
import { CompetitionAPI, ParticipantAPI } from "../services/api";
import ReturnButton from "../components/ReturnButton";

export default function CompetitionDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [competition, setCompetition] = useState<Competition | null>(null);
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const compResponse = await CompetitionAPI.getById(Number(id));
        setCompetition(compResponse.data);

        const partResponse = await ParticipantAPI.getAll(Number(id));
        setParticipants(partResponse.data);
      } catch (err) {
        setError("Failed to fetch competition details");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [id]);

  const handleDeleteCompetition = async () => {
    if (window.confirm("Are you sure you want to delete this competition?")) {
      try {
        await CompetitionAPI.delete(Number(id));
        navigate("/competitions");
      } catch (err) {
        setError("Failed to delete competition");
      }
    }
  };

  const handleRemoveParticipant = async (participantId: number) => {
    if (
      window.confirm(
        "Are you sure you want to remove this participant from the competition?"
      )
    ) {
      try {
        await ParticipantAPI.removeFromCompetition(participantId, Number(id));
        setParticipants(participants.filter((p) => p.id !== participantId));
      } catch (err) {
        setError("Failed to remove participant");
      }
    }
  };

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!competition) return <div className="error">Competition not found</div>;

  return (
    <div className="competition-details">
      <ReturnButton to="/competitions" />
      <div className="header">
        <h1>{competition.name}</h1>
        <div className="actions">
          <Link to={`/competitions/${id}/edit`} className="btn edit">
            Edit
          </Link>
          <button onClick={handleDeleteCompetition} className="btn danger">
            Delete
          </button>
        </div>
      </div>

      <div className="details-grid">
        <div className="detail-item">
          <label>Date:</label>
          <p>{new Date(competition.date).toLocaleDateString()}</p>
        </div>
        <div className="detail-item">
          <label>Location:</label>
          <p>{competition.location}</p>
        </div>
        <div className="detail-item description">
          <label>Description:</label>
          <p>{competition.description}</p>
        </div>
      </div>

      <div className="participants-section">
        <h2>Participants ({participants.length})</h2>
        <Link to={`/participants/new?competition_id=${id}`} className="btn add">
          Add New Participant
        </Link>

        {participants && (
          <div className="participants-list">
            {participants.map((participant) => (
              <div key={participant.id} className="participant-card">
                <div className="participant-info">
                  <h3>{participant.name}</h3>
                  <p>{participant.email}</p>
                </div>
                <div className="participant-actions">
                  <Link
                    to={`/participants/${participant.id}/edit`}
                    className="btn edit"
                  >
                    Edit
                  </Link>
                  <button
                    onClick={() => handleRemoveParticipant(participant.id)}
                    className="btn danger"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
