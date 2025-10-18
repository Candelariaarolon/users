import { useNavigate } from "react-router-dom";

// Hook personalizado para manejar la autenticación
function useAuth() {
  const navigate = useNavigate();

  const getUserIDFromToken = () => {
    try {
      const tokenCookie = document.cookie
        .split("; ")
        .find((row) => row.startsWith("token="));

      if (!tokenCookie) {
        alert("No estás autenticado. Por favor, inicia sesión.");
        navigate("/login");
        return null;
      }

      const token = tokenCookie.split("=")[1];

      // Verificar que el token tenga el formato JWT (3 partes separadas por puntos)
      const tokenParts = token.split(".");
      if (tokenParts.length !== 3) {
        throw new Error("Token no tiene el formato JWT válido");
      }

      const payload = JSON.parse(atob(tokenParts[1]));

      // Verificar diferentes posibles nombres para el user ID
      const userID = payload.jti;

      if (!userID) {
        console.error("No se encontró user_id en el payload del token");
        alert(
          "No se pudo obtener el ID de usuario. Por favor, inicia sesión nuevamente."
        );
        navigate("/login");
        return null;
      }

      const parsedUserID = parseInt(userID);

      if (isNaN(parsedUserID)) {
        console.error("User ID no es un número válido");
        alert("ID de usuario inválido. Por favor, inicia sesión nuevamente.");
        navigate("/login");
        return null;
      }

      return parsedUserID;
    } catch (error) {
      console.error("Error al procesar el token:", error);
      console.error("Stack trace:", error.stack);
      alert("Token inválido. Por favor, inicia sesión nuevamente.");
      navigate("/login");
      return null;
    }
  };
  const id = getUserIDFromToken();

  return id;
}

export default useAuth;
