package clients

import (
	"backend/model"
	"fmt"

	"gorm.io/gorm"
)

var Db *gorm.DB

func GetUserByUsername(username string) (model.UserModel, error) {
	var user model.UserModel
	query := Db.First(&user, "email = ?", username)
	if query.Error != nil {
		return model.UserModel{}, fmt.Errorf("failed to get user by username: %w", query.Error)
	}

	return user, nil
}

// getuserbyid
func GetUserByID(id int) (model.UserModel, error) {
	var user model.UserModel
	// Es equivalente a hacer un select * from users where id = id
	query := Db.Where("id = ?", id).Preload("InscriptionsUser").First(&user)
	if query.Error != nil {
		return model.UserModel{}, fmt.Errorf("failed to get user by id: %w", query.Error)
	}

	if user.ID == 0 {
		err := gorm.ErrRecordNotFound
		return model.UserModel{}, err
	}

	return user, nil
}

func GetUserActivities(UserId int) ([]model.ActivityModel, error) {
	var user model.UserModel

	Db.Preload("InscriptionsUser.Activity").First(&user, UserId)

	var activitiesArray model.Activities
	for _, insc := range user.InscriptionsUser {
		activitiesArray = append(activitiesArray, insc.Activity)
	}

	if len(activitiesArray) == 0 {
		err := gorm.ErrRecordNotFound
		return activitiesArray, err
	}
	return activitiesArray, nil
}
