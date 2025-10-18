package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var request dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	response, err := services.Register(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func VerifyEmail(ctx *gin.Context) {
	var request dto.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	response, err := services.VerifyEmail(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func ResendVerificationCode(ctx *gin.Context) {
	var request dto.ResendCodeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err := services.ResendVerificationCode(request.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

func Login(ctx *gin.Context) {
	var request dto.LoginRequest
	// recibo usuario y contraseña desde el body de la request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// llamar al servicio de login
	// el servicio de login devuelve access token, refresh token, nombre y apellido
	response, err := services.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No se pudo iniciar sesion"})
		return
	}

	// si el login es exitoso, devolver la respuesta con ambos tokens
	ctx.JSON(http.StatusOK, response)
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

func RefreshToken(ctx *gin.Context) {
	var request dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// llamar al servicio de refresh token
	response, err := services.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
