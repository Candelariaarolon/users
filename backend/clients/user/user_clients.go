package clients

import (
	"backend/model"
	"fmt"
	"time"

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

// GetUserByID gets a user by ID
func GetUserByID(id int) (model.UserModel, error) {
	var user model.UserModel
	query := Db.Where("id = ?", id).First(&user)
	if query.Error != nil {
		return model.UserModel{}, fmt.Errorf("failed to get user by id: %w", query.Error)
	}

	if user.ID == 0 {
		err := gorm.ErrRecordNotFound
		return model.UserModel{}, err
	}

	return user, nil
}

// CreateUser creates a new user in the database
func CreateUser(user model.UserModel) (model.UserModel, error) {
	result := Db.Create(&user)
	if result.Error != nil {
		return model.UserModel{}, fmt.Errorf("failed to create user: %w", result.Error)
	}
	return user, nil
}

// UpdateUser updates an existing user
func UpdateUser(user model.UserModel) error {
	result := Db.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

// GetUserByEmail gets a user by email address
func GetUserByEmail(email string) (model.UserModel, error) {
	var user model.UserModel
	query := Db.Where("email = ?", email).First(&user)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return model.UserModel{}, gorm.ErrRecordNotFound
		}
		return model.UserModel{}, fmt.Errorf("failed to get user by email: %w", query.Error)
	}
	return user, nil
}

// VerifyUserEmail verifies user email and clears verification code
func VerifyUserEmail(userID int) error {
	result := Db.Model(&model.UserModel{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_verified":       true,
			"verification_code": nil,
			"code_expires_at":   nil,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to verify user email: %w", result.Error)
	}
	return nil
}

// UpdateVerificationCode updates the verification code and expiration
func UpdateVerificationCode(userID int, code string, expiresAt time.Time) error {
	result := Db.Model(&model.UserModel{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"verification_code": code,
			"code_expires_at":   expiresAt,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update verification code: %w", result.Error)
	}
	return nil
}

// PromoteToAdmin promotes a user to admin status
func PromoteToAdmin(userID int) error {
	result := Db.Model(&model.UserModel{}).
		Where("id = ?", userID).
		Update("is_admin", true)
	if result.Error != nil {
		return fmt.Errorf("failed to promote user to admin: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
