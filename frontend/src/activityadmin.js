import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useState, useRef } from "react";
import "./style_activities.css";

function SaludoUsuario() {
  const userName = localStorage.getItem("userName");
  const userSurname = localStorage.getItem("surname");
  return (
    <div>
      <h2 className="user-welcome">
        Hello {userName} {userSurname}!
      </h2>
    </div>
  );
}

const Activity = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [activity, setActivity] = useState(null);
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
        Authorization: `${token}`, // Este fragmento establece una cabecera HTTP
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
        return res.json().catch((err) => {
          // ✅ AGREGAR manejo de errores
          if (!err.message.includes("Unauthorized")) {
            console.error("Error fetching activity:", err);
          }
        });
      })
      .then((data) => setActivity(data));
  }, [id, navigate]);

  if (!activity) return <p className="text-center mt-5">Loading...</p>;

  return (
    <div className="container my-5">
      <SaludoUsuario />
      <button className="volver-btn" onClick={() => navigate(`/activities`)}>
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
          <p>
            <strong>Quotas Available:</strong> {activity.quotas_available}
          </p>
          <button
            className="btn btn-danger mt-3"
            onClick={() => {
              const token = document.cookie
                .split("; ")
                .find((row) => row.startsWith("token="))
                ?.split("=")[1];
              if (!token) {
                alert("No estás autenticado. Por favor, inicia sesión.");
                navigate("/login");
                return;
              }
              if (
                window.confirm(
                  "¿Estás seguro de que deseas eliminar esta actividad?"
                )
              ) {
                fetch(`http://localhost:8080/activity/${activity.id}`, {
                  method: "DELETE",
                  headers: {
                    "Content-Type": "application/json",
                    Authorization: `${token}`,
                  },
                })
                  .then((res) => {
                    if (res.ok) {
                      alert("Actividad eliminada con éxito");
                      navigate("/users/admin");
                    } else if (res.status === 401 || res.status === 403) {
                      alert("No estás autenticado. Por favor, inicia sesión.");
                      navigate("/login");
                    } else {
                      alert(
                        `Error: ${res.status}. No se pudo eliminar la actividad.`
                      );
                    }
                  })
                  .catch(() => alert("Error al eliminar la actividad"));
              }
            }}
          >
            Eliminar
          </button>
          <button
            className="btn btn-dark mt-3"
            onClick={() => navigate(`/editactivity/${activity.id}`)}
          >
            Editar
          </button>
        </div>
      </div>
    </div>
  );
};

export default Activity;
