import { useState, useEffect } from "react";
import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { ParticipantFormData } from "../types/types";
import { ParticipantAPI, CompetitionAPI } from "../services/api";
import ReturnButton from "../components/ReturnButton";

export default function ParticipantForm() {
  const { id } = useParams();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [formData, setFormData] = useState<ParticipantFormData>({
    name: "",
    email: "",
  });
  const [competitions, setCompetitions] = useState<
    { id: number; name: string }[]
  >([]);
  const [selectedCompetition, setSelectedCompetition] = useState<number>(
    Number(searchParams.get("competition_id")) || 0
  );
  const [registrationDate, setRegistrationDate] = useState<string>(
    new Date().toISOString().split("T")[0]
  );
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch competitions for the dropdown
        const compResponse = await CompetitionAPI.getAll();
        setCompetitions(compResponse.data);

        if (id) {
          const response = await ParticipantAPI.getById(Number(id));
          setFormData(response.data);
        }
      } catch (err) {
        console.error("Error fetching data:", err);
      }
    };
    fetchData();
  }, [id]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Name is required";
    } else if (formData.name.length > 255) {
      newErrors.name = "Name is too long (maximum 255 characters)";
    }

    if (!formData.email.trim()) {
      newErrors.email = "Email is required";
    } else {
      const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
      if (!emailRegex.test(formData.email)) {
        newErrors.email = "Invalid email format";
      }
    }

    if (!selectedCompetition || selectedCompetition <= 0) {
      newErrors.competition = "Competition is required";
    }

    if (!registrationDate) {
      newErrors.registration_date = "Registration date is required";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    setLoading(true);
    try {
      if (id) {
        // Update participant
        await ParticipantAPI.update(Number(id), formData);
        // Add to competition if not already added
        await ParticipantAPI.addToCompetition(Number(id), {
          competition_id: selectedCompetition,
          registration_date: registrationDate,
        });
      } else {
        // Create participant
        const response = await ParticipantAPI.create(formData);
        // Add to competition
        await ParticipantAPI.addToCompetition(response.data.id, {
          competition_id: selectedCompetition,
          registration_date: registrationDate,
        });
      }
      navigate(`/competitions/${selectedCompetition}`);
    } catch (err) {
      console.error("Error saving participant:", err);
      setErrors({ general: "Failed to save participant. Please try again." });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="participant-form">
      <ReturnButton to={`/competitions/${selectedCompetition}`} />
      <h1>{id ? "Edit" : "Add"} Participant</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            className={errors.name ? "error" : ""}
            maxLength={255}
          />
          {errors.name && <span className="error-message">{errors.name}</span>}
        </div>

        <div className="form-group">
          <label>Email:</label>
          <input
            type="email"
            value={formData.email}
            onChange={(e) =>
              setFormData({ ...formData, email: e.target.value })
            }
            className={errors.email ? "error" : ""}
          />
          {errors.email && (
            <span className="error-message">{errors.email}</span>
          )}
        </div>

        <div className="form-group">
          <label>Competition:</label>
          <select
            value={selectedCompetition}
            onChange={(e) => setSelectedCompetition(Number(e.target.value))}
            className={errors.competition ? "error" : ""}
          >
            <option value="">Select a competition</option>
            {competitions.map((comp) => (
              <option key={comp.id} value={comp.id}>
                {comp.name}
              </option>
            ))}
          </select>
          {errors.competition && (
            <span className="error-message">{errors.competition}</span>
          )}
        </div>

        <div className="form-group">
          <label>Registration Date:</label>
          <input
            type="date"
            value={registrationDate}
            onChange={(e) => setRegistrationDate(e.target.value)}
            className={errors.registration_date ? "error" : ""}
          />
          {errors.registration_date && (
            <span className="error-message">{errors.registration_date}</span>
          )}
        </div>

        <div className="form-actions">
          <button type="submit" className="btn primary" disabled={loading}>
            {loading ? "Saving..." : "Save"}
          </button>
          <button
            type="button"
            className="btn secondary"
            onClick={() => navigate(-1)}
          >
            Cancel
          </button>
        </div>

        {errors.general && (
          <div className="error-message">{errors.general}</div>
        )}
      </form>
    </div>
  );
}
