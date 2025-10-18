package services

import (
	"backend/model"
	"fmt"
	"log"
	"time"

	userCLient "backend/clients/user"
	"backend/dto"
	"backend/utils"

	"gorm.io/gorm"
)

func Register(request dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Check if user already exists
	existingUser, err := userCLient.GetUserByEmail(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Error checking existing user:", err)
		return dto.RegisterResponse{}, fmt.Errorf("error checking user existence: %w", err)
	}

	if existingUser.ID != 0 {
		return dto.RegisterResponse{}, fmt.Errorf("user with email %s already exists", request.Email)
	}

	// Hash password
	passwordHash := utils.HashSHA256(request.Password)

	// Generate verification code
	verificationCode, err := utils.GenerateVerificationCode()
	if err != nil {
		log.Println("Error generating verification code:", err)
		return dto.RegisterResponse{}, fmt.Errorf("error generating verification code: %w", err)
	}

	// Create user
	newUser := model.UserModel{
		Email:            request.Email,
		PasswordHash:     passwordHash,
		FirstName:        request.FirstName,
		LastName:         request.LastName,
		IsAdmin:          false,
		IsVerified:       false,
		VerificationCode: verificationCode,
		CodeExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	createdUser, err := userCLient.CreateUser(newUser)
	if err != nil {
		log.Println("Error creating user:", err)
		return dto.RegisterResponse{}, fmt.Errorf("error creating user: %w", err)
	}

	// Send verification email
	err = utils.SendVerificationEmail(createdUser.Email, verificationCode, createdUser.FirstName)
	if err != nil {
		log.Println("Error sending verification email:", err)
		// Don't fail registration if email fails
	}

	return dto.RegisterResponse{
		Message: "User registered successfully. Please check your email for verification code.",
		Email:   createdUser.Email,
	}, nil
}

func VerifyEmail(request dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	// Get user by email
	user, err := userCLient.GetUserByEmail(request.Email)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return dto.VerifyEmailResponse{}, fmt.Errorf("user not found")
	}

	// Check if already verified
	if user.IsVerified {
		return dto.VerifyEmailResponse{}, fmt.Errorf("email already verified")
	}

	// Check if code matches
	if user.VerificationCode != request.Code {
		return dto.VerifyEmailResponse{}, fmt.Errorf("invalid verification code")
	}

	// Check if code expired
	if time.Now().After(user.CodeExpiresAt) {
		return dto.VerifyEmailResponse{}, fmt.Errorf("verification code expired")
	}

	// Verify user
	err = userCLient.VerifyUserEmail(user.ID)
	if err != nil {
		log.Println("Error verifying user email:", err)
		return dto.VerifyEmailResponse{}, fmt.Errorf("error verifying email: %w", err)
	}

	// Send welcome email
	err = utils.SendWelcomeEmail(user.Email, user.FirstName)
	if err != nil {
		log.Println("Error sending welcome email:", err)
		// Don't fail verification if welcome email fails
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, user.IsAdmin)
	if err != nil {
		log.Println("Error generating tokens:", err)
		return dto.VerifyEmailResponse{
			Message: "Email verified successfully",
		}, nil
	}

	return dto.VerifyEmailResponse{
		Message:      "Email verified successfully. You can now log in.",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ResendVerificationCode(email string) error {
	// Get user by email
	user, err := userCLient.GetUserByEmail(email)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return fmt.Errorf("user not found")
	}

	// Check if already verified
	if user.IsVerified {
		return fmt.Errorf("email already verified")
	}

	// Generate new verification code
	verificationCode, err := utils.GenerateVerificationCode()
	if err != nil {
		log.Println("Error generating verification code:", err)
		return fmt.Errorf("error generating verification code: %w", err)
	}

	// Update verification code
	err = userCLient.UpdateVerificationCode(user.ID, verificationCode, time.Now().Add(15*time.Minute))
	if err != nil {
		log.Println("Error updating verification code:", err)
		return fmt.Errorf("error updating verification code: %w", err)
	}

	// Send verification email
	err = utils.SendVerificationEmail(user.Email, verificationCode, user.FirstName)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return fmt.Errorf("error sending verification email: %w", err)
	}

	return nil
}

func Login(username string, password string) (dto.LoginResponse, error) {
	userModel, err := userCLient.GetUserByUsername(username)
	if err != nil {
		log.Println("Error al obtener el usuario por username")
		return dto.LoginResponse{}, fmt.Errorf("failed to get user by user: %w", err)
	}

	// Check if email is verified
	if !userModel.IsVerified {
		log.Println("User email not verified")
		return dto.LoginResponse{}, fmt.Errorf("please verify your email before logging in")
	}

	if utils.HashSHA256(password) != userModel.PasswordHash {
		log.Println("Error al obtener el usuario por password")
		return dto.LoginResponse{}, fmt.Errorf("invalid password")
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(userModel.ID, userModel.IsAdmin)
	if err != nil {
		log.Println("Error al generar los tokens")
		return dto.LoginResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userModel.FirstName,
		Surname:      userModel.LastName,
	}, nil
}

func GetUserByID(id int) (dto.UserDto, error) {
	userModel, err := userCLient.GetUserByID(id)
	if err != nil {
		return dto.UserDto{}, err
	}

	return dto.UserDto{
		ID:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
		IsAdmin:   userModel.IsAdmin,
	}, err
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

// RefreshAccessToken validates a refresh token and generates new access and refresh tokens
func RefreshAccessToken(refreshToken string) (dto.RefreshTokenResponse, error) {
	// Validate refresh token and extract user info
	userID, isAdmin, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Println("Error validating refresh token:", err)
		return dto.RefreshTokenResponse{}, fmt.Errorf("invalid or expired refresh token")
	}

	// Generate new token pair
	newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(userID, isAdmin)
	if err != nil {
		log.Println("Error generating new token pair:", err)
		return dto.RefreshTokenResponse{}, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
