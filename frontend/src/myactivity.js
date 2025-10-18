import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useState, useRef } from "react";
import "./style_activities.css";

const MyActivity = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [activity, setActivity] = useState(null);
  const hasRun = useRef(false);

  const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("token="))
    ?.split("=")[1];

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;
    if (!token) {
      alert("No estás autenticado. Por favor, inicia sesión.");
      navigate("/login");
      return;
    }

    fetch(`http://localhost:8080/activities/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `${token}`,
      },
    })
      .then((res) => {
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
      .then((data) => setActivity(data))
      .catch((err) => {
        if (!err.message.includes("Unauthorized")) {
          console.error("Error fetching activity:", err);
        }
      });
  }, [id, token, navigate]);

  if (!activity) return <p className="text-center mt-5">Cargando...</p>;

  return (
    <div className="container my-5">
      <button
        className="volver-btn"
        onClick={() => navigate(`/users/activities`)}
      >
        Go back
      </button>
      <div className="card-activity">
        <img src={activity.photo} alt="Actividad" />
        <div className="card-activity-body">
          <h2 className="card-activity-title">{activity.name}</h2>
          <p>
            <strong>Description:</strong> {activity.description}
          </p>
          <p>
            <strong>Category:</strong> {activity.category}
          </p>
          <p>
            <strong>Professor:</strong> {activity.profesor_name}
          </p>
          <p>
            <strong>Day:</strong> {activity.day}
          </p>
          <p>
            <strong>Start Hour:</strong> {activity.hour_start}
          </p>
        </div>
      </div>
    </div>
  );
};

export default MyActivity;
