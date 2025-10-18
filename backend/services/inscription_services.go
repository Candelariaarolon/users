package services

import (
	inscriptionClient "backend/clients/inscription"
	"backend/dto"
	"backend/model"
	"time"
)

func Inscription(user_id int, activity_id int) (int, error) {

	activitycheck, err := GetActivityByID(activity_id)

	for _, ins := range activitycheck.InscriptionsActivity {
		if err != nil {
			return 0, err // Return error if activity not found
		}
		if ins.UserID == user_id { // User is already inscribed
			return 1, err // User is already inscribed, no action needed
		}
	}

	if activitycheck.QuotasAvailable <= 0 { // No available quotas
		return 2, err // No available quotas, cannot inscribe
	}

	inscription := dto.InscriptionDto{
		UserID:          user_id,
		ActivityID:      activity_id,
		DateInscription: time.Now().Format("2006-01-02"),
		Active:          true,
	}

	inscriptionModel := model.InscriptionModel{
		UserID:          inscription.UserID,
		ActivityID:      inscription.ActivityID,
		DateInscription: inscription.DateInscription,
		Active:          inscription.Active,
	}

	err1 := inscriptionClient.CreateInscription(inscriptionModel)
	if err1 != nil {
		return 0, err1
	}
	return 0, nil
}
