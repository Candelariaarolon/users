package controllers

import (
	"backend/dto"
	"backend/services" //importo modulo propio
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin" //importo un link
)

func GetActivityByID(ctx *gin.Context) {

	activityIDString := ctx.Param("id")                  //me devuelve algo que quiero de la url, le debo poner el mismo nombre
	activityIDInt, err := strconv.Atoi(activityIDString) //convierto el dato en int, el err representa el error que devuelve el cambio de stoi, si no funciona, es una variable
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}
	activity, err1 := services.GetActivityByID(activityIDInt) //le paso el id de la url
	if err1 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, activity) //devuelvo un JSON
}

func GetFilteredActivities(ctx *gin.Context) { //context contiene la info de la request y la response, yo la edito en esta funcion con lo que quiero

	category := ctx.Query("category")
	name := ctx.Query("name")
	description := ctx.Query("description")
	day := ctx.Query("day")
	hour_start := ctx.Query("hour_start")
	professor_name := ctx.Query("professor_name")

	activities, err := services.GetFilteredActivities(category, name, description, day, hour_start, professor_name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No hay actividades"})
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func CreateActivity(ctx *gin.Context) {
	//obtengo los datos del body de la request

	var activity dto.ActivityRequestDto
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		fmt.Printf("Error ShouldBindJSON: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos",
		})
		return
	}

	if !VerifyActivity(activity, ctx) {
		return //si no verifico la actividad, no sigo
	}

	err := services.CreateActivity(activity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la actividad"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Actividad creada con exito"})
}

func ModifyActivity(ctx *gin.Context) {

	activityIDString := ctx.Param("id")                  //me devuelve algo que quiero de la url, le debo poner el mismo nombre
	activityIDInt, err := strconv.Atoi(activityIDString) //convierto el dato en int, el err representa el error que devuelve el cambio de stoi, si no funciona, es una variable
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}
	activity, err1 := services.GetActivityByID(activityIDInt) //le paso el id de la url
	if err1 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
		return
	}

	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		//Gin lee el JSON enviado en el body del request
		return
	}
	activity.ID = activityIDInt

	if !VerifyActivity(ToRequestDto(activity), ctx) {
		return //si no verifico la actividad, no sigo
	}

	err2 := services.ModifyActivity(activity) //le paso el id de la url y los datos de la actividad
	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Actividad modificada con exito"})

}

func DeleteActivity(ctx *gin.Context) {
	activityIDString := ctx.Param("id")                  //me devuelve algo que quiero de la url, le debo poner el mismo nombre
	activityIDInt, err := strconv.Atoi(activityIDString) //convierto el dato en int, el err representa el error que devuelve el cambio de stoi, si no funciona, es una variable
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id invalido"})
		return
	}
	err1 := services.DeleteActivity(activityIDInt) //le paso el id de la url
	if err1 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Actividad eliminada con exito"})
}

// Función para convertir ActivityDto a ActivityRequestDto
func ToRequestDto(a dto.ActivityDto) dto.ActivityRequestDto {
	return dto.ActivityRequestDto{
		Category:     a.Category,
		Name:         a.Name,
		Description:  a.Description,
		ProfesorName: a.ProfesorName,
		Quotas:       a.Quotas,
		Day:          a.Day,
		HourStart:    a.HourStart,
		Active:       a.Active,
		Photo:        a.Photo,
	}
}

func VerifyActivity(activity dto.ActivityRequestDto, ctx1 *gin.Context) bool {

	if activity.Name == "" || activity.Description == "" || activity.Category == "" ||
		activity.ProfesorName == "" || activity.Day == "" || activity.HourStart == "" || activity.Photo == "" {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Hay campos obligatorios vacíos"})
		return false
	}

	//verificar que el nombre del profesor no tenga números
	if strings.ContainsAny(activity.ProfesorName, "0123456789") {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "El nombre del profesor no puede contener números"})
		return false
	}

	timeParts := strings.Split(activity.HourStart, ":")
	if len(timeParts) != 3 {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Formato de hora inválido. Debe ser HH:MM:SS"})
		return false
	}

	hours := timeParts[0]
	minutes := timeParts[1]
	seconds := timeParts[2]

	// Validar que cada parte sea un número válido
	hoursInt, err1 := strconv.Atoi(hours)
	if err1 != nil || hoursInt < 0 || hoursInt > 23 {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Horas inválidas (0-23)"})
		return false
	}

	minutesInt, err2 := strconv.Atoi(minutes)
	if err2 != nil || minutesInt < 0 || minutesInt > 59 {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Minutos inválidos (0-59)"})
		return false
	}

	secondsInt, err3 := strconv.Atoi(seconds)
	if err3 != nil || secondsInt < 0 || secondsInt > 59 {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Segundos inválidos (0-59)"})
		return false
	}

	if activity.Quotas <= 0 {
		ctx1.JSON(http.StatusBadRequest, gin.H{"error": "Los cupos deben ser mayores a 0"})
		return false
	}

	return true
}
