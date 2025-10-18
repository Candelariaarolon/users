package services

import (
	activityClient "backend/clients/activity"
	"backend/dto"
	"backend/model"
	"log"
)

func GetActivityByID(id int) (dto.ActivityDto, error) {
	activity, err := activityClient.GetActivityByID(id)
	var activityDto dto.ActivityDto
	var insDto dto.Inscriptions

	for _, ins := range activity.InscriptionsActivity {

		insDto = append(insDto, dto.InscriptionDto{
			ID:              ins.ID,
			UserID:          ins.UserID,
			ActivityID:      ins.ActivityID,
			DateInscription: ins.DateInscription,
			Active:          ins.Active,
		})
	}

	activityDto.ID = activity.ID
	activityDto.Category = activity.Category
	activityDto.Name = activity.Name
	activityDto.Description = activity.Description
	activityDto.ProfesorName = activity.ProfesorName
	activityDto.Quotas = activity.Quotas
	activityDto.QuotasAvailable = activity.Quotas - len(activity.InscriptionsActivity)
	activityDto.Day = activity.Day
	activityDto.HourStart = activity.HourStart
	activityDto.InscriptionsActivity = insDto
	activityDto.Active = activity.Active
	activityDto.Photo = activity.Photo
	return activityDto, err
}

func GetFilteredActivities(category string, name string, description string, day string, hour_start string, professor_name string) (dto.ActivitiesDto, error) {
	var activitiesDto dto.ActivitiesDto

	activities, err := activityClient.GetFilteredActivities(category, name, description, day, hour_start, professor_name)

	for _, activity := range activities {
		var insDto dto.Inscriptions
		for _, ins := range activity.InscriptionsActivity {
			insDto = append(insDto, dto.InscriptionDto{
				ID:              ins.ID,
				DateInscription: ins.DateInscription,
				Active:          ins.Active,
			})
		}
		activityDto := dto.ActivityDto{
			ID:                   activity.ID,
			Category:             activity.Category,
			Name:                 activity.Name,
			Description:          activity.Description,
			ProfesorName:         activity.ProfesorName,
			Quotas:               activity.Quotas,
			QuotasAvailable:      activity.Quotas - len(activity.InscriptionsActivity),
			Day:                  activity.Day,
			HourStart:            activity.HourStart,
			InscriptionsActivity: insDto,
			Active:               activity.Active,
			Photo:                activity.Photo,
		}
		activitiesDto = append(activitiesDto, activityDto)
	}

	return activitiesDto, err
}

func CreateActivity(activity dto.ActivityRequestDto) error {
	activityModel := model.ActivityModel{
		Category:     activity.Category,
		Name:         activity.Name,
		Description:  activity.Description,
		ProfesorName: activity.ProfesorName,
		Quotas:       activity.Quotas,
		Day:          activity.Day,
		HourStart:    activity.HourStart,
		Active:       activity.Active,
		Photo:        activity.Photo,
	}

	_, err := activityClient.CreateActivity(activityModel)

	if err != nil {
		log.Printf("Error creating activity: %v", err)
		return err
	}

	return nil
}

func ModifyActivity(activity dto.ActivityDto) error {
	activityModel := model.ActivityModel{
		ID:           activity.ID,
		Category:     activity.Category,
		Name:         activity.Name,
		Description:  activity.Description,
		ProfesorName: activity.ProfesorName,
		Quotas:       activity.Quotas,
		Day:          activity.Day,
		HourStart:    activity.HourStart,
		Active:       activity.Active,
		Photo:        activity.Photo,
	}

	err := activityClient.ModifyActivity(activityModel)
	if err != nil {
		return err
	}

	return nil
}

func DeleteActivity(id int) error {
	err := activityClient.DeleteActivity(id)

	if err != nil {
		return err
	}

	return nil
}
