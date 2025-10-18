package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var request dto.LoginRequest
	// recibo usuario y contraseña desde el body de la request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// llamar al servicio de login
	// el servicio de login devuelve el id del usuario y el token
	token, name, surname, err := services.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No se pudo iniciar sesion"})
		return
	}
	// si el login es exitoso, devolver el id del usuario y el token
	// el token es un string que se genera al momento de hacer login
	ctx.JSON(http.StatusOK, dto.LoginResponse{
		Token:   token,
		Name:    name,
		Surname: surname,
	})
}

func GetUserByID(ctx *gin.Context) {
	// recibo el id del usuario desde el path de la request
	userID := ctx.Param("id")
	// hago string a int
	userIDInt, err1 := strconv.Atoi(userID)
	// llamar al servicio de get user by id

	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(userIDInt)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// si el usuario existe, devolver el usuario
	ctx.JSON(http.StatusOK, user)
}

func GetUserActivities(ctx *gin.Context) {
	// recibo el id del usuario desde el path de la request
	userID := ctx.Param("id")
	// hago string a int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// llamar al servicio de get user activities
	activities, err2 := services.GetUserActivities(userIDInt)
	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "El usuario no tiene actividades"})
		return
	}
	// si el usuario existe, devolver las actividades
	ctx.JSON(http.StatusOK, activities)
}

func VerifyToken(ctx *gin.Context) {
	// recibo el token desde el header de la request
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		ctx.Abort()
		return
	}

	// llamar al servicio de verify token
	err := services.VerifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
}

func VerifyAdminToken(ctx *gin.Context) {
	// recibo el token desde el header de la request
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		ctx.Abort()
		return
	}

	// llamar al servicio de verify admin token
	err := services.VerifyAdminToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
}
