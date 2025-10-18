package dto

type InscriptionDto struct {
	ID              int    `json:"id"`               //PK
	UserID          int    `json:"user_id"`          //FK
	ActivityID      int    `json:"activity_id"`      //FK
	DateInscription string `json:"date_inscription"` //Date of inscription
	Active          bool   `json:"active"`           //Active status
}

type Inscriptions []InscriptionDto

type InscriptionRequest struct {
	UserID     int `json:"user_id"`
	ActivityID int `json:"activity_id"`
}
