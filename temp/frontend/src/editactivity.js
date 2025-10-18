import { useState, useEffect, useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import "./editactivity.css";

const EditActivity = () => {
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
  const [original, setOriginal] = useState({});
  const navigate = useNavigate();
  const { id } = useParams();
  const [hour, setHour] = useState("");
  const [minute, setMinute] = useState("");
  const [second, setSecond] = useState("");
  const hasRun = useRef(false);

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;

    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("token="))
      ?.split("=")[1];

    if (!token) {
      alert("You are not authenticated. Please log in.");
      navigate("/login");
      return;
    }

    fetch(`http://localhost:8080/activities/${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `${token}`,
      },
    })
      .then((res) => {
        // ← AGREGAR validación de respuesta
        if (res.status === 401) {
          alert("Your session has expired. Please log in again.");
          navigate("/login");
          throw new Error("Unauthorized");
        }
        if (!res.ok) {
          throw new Error("Failed to fetch activity");
        }
        return res.json();
      })
      .then((data) => {
        setOriginal(data);

        setForm({
          name: data.name || "",
          category: data.category || "",
          description: data.description || "",
          profesor_name: data.profesor_name || "",
          quotas: data.quotas?.toString() || "",
          day: data.day || "",
          hour_start: data.hour_start || "",
          active: data.active ?? true,
          photo: data.photo || "",
        });

        if (data.hour_start) {
          const [h, m, s] = data.hour_start.split(":");
          setHour(h || "");
          setMinute(m || "");
          setSecond(s || "");
        }
      })
      .catch((err) => {
        // ← AGREGAR manejo de errores
        if (!err.message.includes("Unauthorized")) {
          console.error("Error fetching activity:", err);
          alert("Error loading activity data");
        }
      });
  }, [id, navigate]);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleHourChange = (e) => {
    const value = e.target.value.replace(/\D/, "");
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

    if (!token) {
      alert("You are not authenticated. Please log in.");
      navigate("/login");
      return;
    }

    const hour_start = `${hour.padStart(2, "0")}:${minute.padStart(
      2,
      "0"
    )}:${second.padStart(2, "0")}`;

    const formToSend = {
      name: form.name || "",
      category: form.category || "",
      description: form.description || "",
      profesor_name: form.profesor_name || "",
      quotas: form.quotas ? Number(form.quotas) : 0,
      day: form.day || "",
      photo: form.photo || "",
      active: form.active,
      hour_start:
        hour !== "" || minute !== "" || second !== "" ? hour_start : "",
    };
    Object.keys(form).forEach((key) => {
      if (form[key] !== "" && form[key] !== null && form[key] !== undefined) {
        formToSend[key] = key === "quotas" ? Number(form[key]) : form[key];
      }
    });

    formToSend.hour_start =
      hour !== "" || minute !== "" || second !== ""
        ? hour_start
        : original.hour_start;

    fetch(`http://localhost:8080/activity/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `${token}`,
      },
      body: JSON.stringify(formToSend),
    })
      .then((res) => {
        if (res.ok) {
          alert("Actividad editada con éxito");
          navigate("/users/admin");
        } else if (res.status === 401 || res.status === 403) {
          alert("No estás autenticado. Por favor, inicia sesión.");
          navigate("/login");
        } else {
          res.json().then((data) => {
            alert(`Error: ${data.error}. No se pudo crear la actividad.`);
          });
        }
      })
      .catch(() => alert("Error al editar la actividad"));
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
          Editar Actividad
        </button>
      </form>
    </div>
  );
};

export default EditActivity;
