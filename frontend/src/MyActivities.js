import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import "./style_myactivities.css";
import useAuth from "./authToken";

function SaludoUsuario() {
  const userFirstName = localStorage.getItem("userName");
  const userLastName = localStorage.getItem("surname");

  return (
    <div>
      <h2 className="user-welcome">
        Hello {userFirstName} {userLastName}!
      </h2>
    </div>
  );
}

const MyActivities = () => {
  const userID = useAuth();
  const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("token="))
    .split("=")[1];
  const [myActivities, setMyActivities] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const hasRun = useRef(false);

  // Cargar actividades cuando tenemos el userID
  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;
    if (userID) {
      fetch(`http://localhost:8080/users/${userID}/activities`, {
        headers: {
          Authorization: `${token}`,
        },
      })
        .then((res) => {
          if (res.status === 401) {
            alert("No estás autenticado. Por favor, inicia sesión.");
            navigate("/login");
            return null;
          }
          if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
          }
          return res.json();
        })
        .then((data) => {
          setMyActivities(data);
          setLoading(false);
        })
        .catch((err) => {
          console.error("Error fetching user activities:", err);
          setLoading(false);
        });
    }
  }, [userID, token, navigate]);

  if (!userID) {
    return <div>Redirigiendo...</div>;
  }

  return (
    <div>
      <button className="volver-btn" onClick={() => navigate(`/activities`)}>
        Go back
      </button>

      <SaludoUsuario userID={userID} />

      <div className="my-activities-wrapper">
        <h2 className="title">My Activities</h2>
        <div className="my-activities-container">
          {loading ? (
            <p>Cargando actividades...</p>
          ) : (
            <div className="activities-grid">
              {myActivities && myActivities.length > 0 ? (
                myActivities.map((activity) => (
                  <div className="activity-card" key={activity.id}>
                    <h5>{activity.name}</h5>
                    <p>{activity.description}</p>
                    <button
                      className="ver-detalles-btn"
                      onClick={() => navigate(`/myactivities/${activity.id}`)}
                    >
                      view details
                    </button>
                  </div>
                ))
              ) : (
                <p>No tienes actividades registradas.</p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MyActivities;
