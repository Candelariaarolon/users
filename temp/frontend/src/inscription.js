import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import useAuth from "./authToken";

const Inscription = () => {
  const { id: activityId } = useParams();
  const navigate = useNavigate();
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const hasRun = useRef(false);
  const userID = useAuth();

  useEffect(() => {
    if (hasRun.current) return;
    hasRun.current = true;

    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("token="))
      ?.split("=")[1];

    if (!hasSubmitted) {
      fetch("http://localhost:8080/users/inscription", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${token}`,
        },
        body: JSON.stringify({
          user_id: parseInt(userID),
          activity_id: parseInt(activityId),
        }),
      })
        .then((response) => {
          if (response.status === 401) {
            alert("No estás autenticado. Por favor, inicia sesión.");
            navigate("/login");
            return null;
          }
          if (response.status === 403) {
            alert("No es posible inscribirse.");
            navigate("/activities");
            return null;
          }
          if (!response.ok) throw new Error("Error en la inscripción");
          return response.json();
        })
        .then((data) => {
          if (data) {
            alert("Inscripción exitosa!");
            setHasSubmitted(true);
            navigate("/activities");
          }
        })
        .catch((err) => {
          console.error(err);
          alert("Hubo un error al inscribirse.");
        });
    }
  }, [activityId, hasSubmitted, navigate, userID]);
};
export default Inscription;
