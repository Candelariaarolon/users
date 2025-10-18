package model

type InscriptionModel struct {
	ID              int       `gorm:"primaryKey;autoIncrement"` //PK
	User            UserModel `gorm:"foreignKey:UserID"`        //FK
	UserID          int
	DateInscription string        `gorm:"type:date;not null"`    //Date of inscription
	Active          bool          `gorm:"default:true"`          //Active status
	Activity        ActivityModel `gorm:"foreignKey:ActivityID"` //FK
	ActivityID      int
}

type Inscriptions []InscriptionModel
