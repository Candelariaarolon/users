package model

import "time"

type UserModel struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`          //PK
	Email            string    `gorm:"unique;not null;type:varchar(100)"` //Unique email
	PasswordHash     string    `gorm:"longtext"`                          //Password Hash
	FirstName        string    `gorm:"type:varchar(100);not null"`
	LastName         string    `gorm:"type:varchar(100);not null"`
	IsAdmin          bool      `gorm:"default:false"`        //Admin
	IsVerified       bool      `gorm:"default:false"`        //Email verified
	CreatedAt        time.Time `gorm:"autoCreateTime"`       //Creation timestamp
	VerificationCode string    `gorm:"type:varchar(6);null"` //6-digit verification code
	CodeExpiresAt    time.Time `gorm:"null"`                 //Code expiration time
}

type VerificationToken struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    int       `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(6);not null"` //6-digit code
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
