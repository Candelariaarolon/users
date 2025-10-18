package clients

import (
	"backend/model"
	"backend/utils/errors"

	"gorm.io/gorm"
)

var Db *gorm.DB

func CreateInscription(inscription model.InscriptionModel) errors.ApiError {
	if err := Db.Create(&inscription).Error; err != nil {
		return errors.NewInternalServerApiError("error creating inscription", err)
	}
	return nil
}
