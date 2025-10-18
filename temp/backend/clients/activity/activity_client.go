package clients

import (
	"backend/model"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Db *gorm.DB

func GetActivityByID(id int) (model.ActivityModel, error) {
	var activity model.ActivityModel
	var err error
	Db.Preload("InscriptionsActivity").First(&activity, id)

	if activity.ID == 0 {
		log.Error("Activity not found")
		err = gorm.ErrRecordNotFound
		return activity, err
	}

	log.Debugf("Activity Found: %+v", activity)

	return activity, err //Devuelve un DAO (Data Access Object) de la actividad
}

func GetFilteredActivities(category string, name string, description string, day string, hour_start string, professor_name string) (model.Activities, error) {
	var activities model.Activities
	var err error

	query := Db // comenzamos con el DB base

	if category != "" {
		query = query.Where("category LIKE ?", "%"+category+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}
	if day != "" {
		query = query.Where("day LIKE ?", "%"+day+"%")
	}
	if hour_start != "" {
		query = query.Where("hour_start LIKE ?", "%"+hour_start+"%")
	}
	if professor_name != "" {
		query = query.Where("profesor_name LIKE ?", "%"+professor_name+"%")
	}

	query.Find(&activities)

	if len(activities) == 0 {
		log.Error("No activities found")
		err = gorm.ErrRecordNotFound
		return activities, err
	}
	log.Debugf("Filtered Activities Found: %+v", activities)

	return activities, err
}

func CreateActivity(activity model.ActivityModel) (model.ActivityModel, error) {
	result := Db.Create(&activity)

	if result.Error != nil {
		log.Error("Error creating activity: ", result.Error)
		return model.ActivityModel{}, result.Error
	}

	log.Debugf("Activity Created: %+v", activity)
	return activity, nil
}

func ModifyActivity(activity model.ActivityModel) error {
	var existingActivity model.ActivityModel
	if err := Db.First(&existingActivity, activity.ID).Error; err != nil {
		log.Error("Error finding activity to modify: ", err)
		return err
	}
	if err := Db.Model(&existingActivity).Updates(activity).Error; err != nil {
		// La función Updates de GORM se utiliza para actualizar múltiples campos de un registro en la base de datos.
		// .Updates(activity) toma los campos no nulos del struct activity y los copia en la base de datos
		log.Error("Error updating activity: ", err)
		return err
	}
	if err := Db.First(&existingActivity, activity.ID).Error; err != nil {
		log.Error("Error fetching updated activity: ", err)
		return err
	}

	log.Debugf("Activity Updated: %+v", existingActivity)
	return nil
}

func DeleteActivity(id int) error {
	// Primero eliminamos inscripciones
	if err := Db.Where("activity_id = ?", id).Delete(&model.InscriptionModel{}).Error; err != nil {
		log.Error("Error deleting inscriptions: ", err)
		return err
	}

	// Luego eliminamos la actividad
	var activity model.ActivityModel
	if err := Db.First(&activity, id).Error; err != nil {
		log.Error("Error finding activity to delete: ", err)
		return err
	}

	if err := Db.Delete(&activity).Error; err != nil {
		log.Error("Error deleting activity: ", err)
		return err
	}

	log.Debugf("Activity Deleted: %d", id)
	return nil
}
