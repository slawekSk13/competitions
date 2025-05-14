import { useEffect, useState } from "react";
import { CompetitionAPI } from "../services/api";
import { Competition } from "../services/api";
import { Link } from "react-router-dom";

export default function CompetitionsList() {
  const [competitions, setCompetitions] = useState<Competition[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCompetitions = async () => {
      try {
        const response = await CompetitionAPI.getAll();
        setCompetitions(response.data);
      } catch (error) {
        console.error("Error fetching competitions:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchCompetitions();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div className="container">
      <h1>Competitions</h1>
      <Link to="/competitions/new" className="btn">
        Create New Competition
      </Link>
      <div className="competition-list">
        {competitions.map((competition) => (
          <div key={competition.id} className="competition-card">
            <h2>{competition.name}</h2>
            <p>{competition.description}</p>
            <div className="actions">
              <Link to={`/competitions/${competition.id}`} className="btn">
                View
              </Link>
              <Link to={`/competitions/${competition.id}/edit`} className="btn">
                Edit
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
