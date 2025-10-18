package services

import (
	userCLient "backend/clients/user"
	"backend/dto"
	"backend/utils"
	"fmt"
	"log"
)

func Login(username string, password string) (string, string, string, error) {
	userModel, err := userCLient.GetUserByUsername(username)
	if err != nil {
		log.Println("Error al obtener el usuario por username")
		return "", "", "", fmt.Errorf("failed to get user by user: %w", err)

	}
	if utils.HashSHA256(password) != userModel.PasswordHash {
		log.Println("Error al obtener el usuario por password")
		return "", "", "", fmt.Errorf("invalid password")

	}

	token, err := utils.GenerateJWT(userModel.ID, userModel.IsAdmin)
	if err != nil {
		log.Println("Error al generar el token")
		return "", "", "", fmt.Errorf("failed to generate token: %w", err)

	}

	name := userModel.FirstName
	surname := userModel.LastName
	return token, name, surname, nil
}

func GetUserByID(id int) (dto.UserDto, error) {
	userModel, err := userCLient.GetUserByID(id)
	if err != nil {
		return dto.UserDto{}, err
	}

	var insDto dto.Inscriptions
	for _, ins := range userModel.InscriptionsUser {
		insDto = append(insDto, dto.InscriptionDto{
			ID:              ins.ID,
			UserID:          ins.UserID,
			ActivityID:      ins.ActivityID,
			DateInscription: ins.DateInscription,
			Active:          ins.Active,
		})

	}
	return dto.UserDto{
		ID:               userModel.ID,
		FirstName:        userModel.FirstName,
		LastName:         userModel.LastName,
		Email:            userModel.Email,
		InscriptionsUser: insDto,
	}, err
}

func GetUserActivities(id int) (dto.ActivitiesDto, error) {
	activityModel, err := userCLient.GetUserActivities(id)
	if err != nil {
		return dto.ActivitiesDto{}, err
	}

	var actDto dto.ActivitiesDto
	for _, act := range activityModel {
		actDto = append(actDto, dto.ActivityDto{
			ID:           act.ID,
			Name:         act.Name,
			ProfesorName: act.ProfesorName,
			Description:  act.Description,
			Category:     act.Category,
			Day:          act.Day,
			HourStart:    act.HourStart,
			Active:       act.Active,
			Photo:        act.Photo,
		})
	}
	return actDto, nil
}

func VerifyToken(token string) error {
	err := utils.ValidateJWT(token)
	if err != nil {
		log.Println("Error al verificar el token")
		return fmt.Errorf("failed to verify token: %w", err)
	}
	return nil
}

func VerifyAdminToken(token string) error {
	err := utils.ValidateAdminJWT(token)
	if err != nil {
		log.Println("Error al verificar el token de admin")
		return fmt.Errorf("failed to verify admin token: %w", err)
	}
	return nil
}
