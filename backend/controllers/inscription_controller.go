package controllers

import (
	"backend/dto"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Inscription(ctx *gin.Context) {

	var request dto.InscriptionRequest
	//recibo el body de la request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//llamo al servicio de inscription
	code, err := services.Inscription(request.UserID, request.ActivityID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Inscription failed"})
		return
	}

	if code == 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User already inscribed"})
		return
	}

	if code == 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No available quotas"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Inscription successful"})
}
