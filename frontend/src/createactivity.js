import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./editactivity.css";

const CreateActivity = () => {
  const [form, setForm] = useState({
    name: "",
    category: "",
    description: "",
    profesor_name: "",
    quotas: "",
    day: "",
    hour_start: "",
    active: true,
    photo: "",
  });
  const navigate = useNavigate();
  const [hour, setHour] = useState("");
  const [minute, setMinute] = useState("");
  const [second, setSecond] = useState("");

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleHourChange = (e) => {
    const value = e.target.value; // Solo números
    if (e.target.name === "hour") setHour(value.slice(0, 2));
    if (e.target.name === "minute") setMinute(value.slice(0, 2));
    if (e.target.name === "second") setSecond(value.slice(0, 2));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("token="))
      ?.split("=")[1];
    const formToSend = { ...form, quotas: Number(form.quotas) };
    formToSend.hour_start = `${hour}:${minute}:${second}`;

    fetch("http://localhost:8080/activity", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `${token}`,
      },
      body: JSON.stringify(formToSend),
      //Así, quotas se convierte a número y el backend puede procesarlo correctamente.
      //Si envías form directamente, quotas sigue siendo string y el backend responde con 400 Bad Request.
    })
      .then((res) => {
        if (res.ok) {
          alert("Actividad creada con éxito");
          navigate("/activities");
        } else if (res.status === 401 || res.status === 403) {
          alert("No estás autenticado. Por favor, inicia sesión.");
          navigate("/login");
        } else {
          res.json().then((data) => {
            alert(`Error: ${data.error}. No se pudo crear la actividad.`);
          });
        }
      })
      .catch(() => alert("Error al crear la actividad"));
  };

  return (
    <div>
      <button className="volver-btn" onClick={() => navigate(`/users/admin`)}>
        Go back
      </button>
      <form onSubmit={handleSubmit} className="edit-activity-form">
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="name">
            Name:
          </label>
          <input
            id="name"
            name="name"
            className="edit-activity-input"
            placeholder="Name"
            value={form.name}
            onChange={handleChange}
          />
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="category">
            Category:
          </label>
          <input
            id="category"
            name="category"
            className="edit-activity-input"
            placeholder="Category"
            value={form.category}
            onChange={handleChange}
          />
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="description">
            Description:
          </label>
          <input
            id="description"
            name="description"
            className="edit-activity-input"
            placeholder="Description"
            value={form.description}
            onChange={handleChange}
          />
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="profesor_name">
            Professor Name:
          </label>
          <input
            id="profesor_name"
            name="profesor_name"
            className="edit-activity-input"
            placeholder="Professor Name"
            value={form.profesor_name}
            onChange={handleChange}
          />
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="quotas">
            Quotas:
          </label>
          <input
            id="quotas"
            name="quotas"
            type="number"
            className="edit-activity-input"
            placeholder="Quotas"
            value={form.quotas}
            onChange={handleChange}
          />
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="day">
            Day:
          </label>
          <select
            id="day"
            name="day"
            className="edit-activity-select"
            value={form.day}
            onChange={handleChange}
            required
          >
            <option value="" disabled>
              Select Day
            </option>
            <option value="Monday">Monday</option>
            <option value="Tuesday">Tuesday</option>
            <option value="Wednesday">Wednesday</option>
            <option value="Thursday">Thursday</option>
            <option value="Friday">Friday</option>
            <option value="Saturday">Saturday</option>
            <option value="Sunday">Sunday</option>
          </select>
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label">Hour Start:</label>
          <div className="edit-activity-hour-group">
            <input
              name="hour"
              type="text"
              placeholder="HH"
              value={hour}
              onChange={handleHourChange}
              maxLength={2}
              className="edit-activity-hour-input"
            />
            :
            <input
              name="minute"
              type="text"
              placeholder="MM"
              value={minute}
              onChange={handleHourChange}
              maxLength={2}
              className="edit-activity-hour-input"
            />
            :
            <input
              name="second"
              type="text"
              placeholder="SS"
              value={second}
              onChange={handleHourChange}
              maxLength={2}
              className="edit-activity-hour-input"
            />
          </div>
        </div>
        <div className="edit-activity-row">
          <label className="edit-activity-label" htmlFor="photo">
            Photo URL:
          </label>
          <input
            id="photo"
            name="photo"
            className="edit-activity-input"
            placeholder="Photo URL"
            value={form.photo}
            onChange={handleChange}
          />
        </div>
        <button
          type="submit"
          className="edit-activity-input"
          style={{
            fontWeight: "bold",
            background: "#db4f0e",
            color: "#fff",
            cursor: "pointer",
          }}
        >
          Crear Actividad
        </button>
      </form>
    </div>
  );
};

export default CreateActivity;
