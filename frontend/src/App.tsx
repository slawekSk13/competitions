import { Routes, Route } from "react-router-dom";
import CompetitionsList from "./pages/CompetitionsList";
import CompetitionDetails from "./pages/CompetitionDetails";
import CompetitionForm from "./pages/CompetitionForm";
import ParticipantForm from "./pages/ParticipantForm";
import "./App.css";

function App() {
  return (
    <Routes>
      <Route path="/" element={<CompetitionsList />} />
      <Route path="/competitions" element={<CompetitionsList />} />
      <Route path="/competitions/new" element={<CompetitionForm />} />
      <Route path="/competitions/:id" element={<CompetitionDetails />} />
      <Route path="/competitions/:id/edit" element={<CompetitionForm />} />
      <Route path="/participants/new" element={<ParticipantForm />} />
      <Route path="/participants/:id/edit" element={<ParticipantForm />} />
    </Routes>
  );
}

export default App;
