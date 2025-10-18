package dto

type ActivityDto struct {
	ID                   int          `json:"id"`       //PK
	Category             string       `json:"category"` //FK
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	ProfesorName         string       `json:"profesor_name"`
	Quotas               int          `json:"quotas"` //Cupos
	Day                  string       `json:"day"`    //FK
	HourStart            string       `json:"hour_start"`
	InscriptionsActivity Inscriptions `json:"inscriptions"`     //FK
	QuotasAvailable      int          `json:"quotas_available"` //Cupos disponibles
	Active               bool         `json:"active"`
	Photo                string       `json:"photo"`
}

type ActivityRequestDto struct {
	Category     string `json:"category"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProfesorName string `json:"profesor_name"`
	Quotas       int    `json:"quotas"`
	Day          string `json:"day"`
	HourStart    string `json:"hour_start"`
	Active       bool   `json:"active"`
	Photo        string `json:"photo"`
}

type ActivitiesDto []ActivityDto
