import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { CompetitionAPI, Competition } from "../services/api";
import ReturnButton from "../components/ReturnButton";

export default function CompetitionForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [formData, setFormData] = useState<Competition>({
    name: "",
    description: "",
    date: "",
    location: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (id) {
      const fetchCompetition = async () => {
        try {
          const response = await CompetitionAPI.getById(Number(id));
          setFormData(response.data);
        } catch (error) {
          console.error("Error fetching competition:", error);
        }
      };
      fetchCompetition();
    }
  }, [id]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Name is required";
    } else if (formData.name.length > 255) {
      newErrors.name = "Name is too long (maximum 255 characters)";
    }

    if (!formData.date) {
      newErrors.date = "Date is required";
    }

    if (!formData.location.trim()) {
      newErrors.location = "Location is required";
    } else if (formData.location.length > 255) {
      newErrors.location = "Location is too long (maximum 255 characters)";
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
        await CompetitionAPI.update(Number(id), formData);
      } else {
        await CompetitionAPI.create(formData);
      }
      navigate("/competitions");
    } catch (error) {
      console.error("Error saving competition:", error);
      setErrors({ general: "Failed to save competition. Please try again." });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="competition-form">
      <ReturnButton to="/competitions" />
      <h1>{id ? "Edit" : "Create"} Competition</h1>
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
          <label>Description:</label>
          <textarea
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
          />
        </div>

        <div className="form-group">
          <label>Date:</label>
          <input
            type="date"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
            className={errors.date ? "error" : ""}
          />
          {errors.date && <span className="error-message">{errors.date}</span>}
        </div>

        <div className="form-group">
          <label>Location:</label>
          <input
            type="text"
            value={formData.location}
            onChange={(e) =>
              setFormData({ ...formData, location: e.target.value })
            }
            className={errors.location ? "error" : ""}
            maxLength={255}
          />
          {errors.location && (
            <span className="error-message">{errors.location}</span>
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
