package app

import (
	"backend/controllers"
	"time"

	"github.com/gin-contrib/cors"
)

func mapUrls() {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, //almacena la configuracion de CORS por 12 horas
	}))

	router.GET("/activities/:id", controllers.VerifyToken, controllers.GetActivityByID)   //trae una actividad por id
	router.GET("/activities", controllers.VerifyToken, controllers.GetFilteredActivities) //trae todas las actividades o filtradas por categoria, nombre o descripcion
	router.POST("/users/login", controllers.Login)
	router.POST("/users/inscription", controllers.VerifyToken, controllers.Inscription)
	router.GET("/users/:id/activities", controllers.VerifyToken, controllers.GetUserActivities) //trae todas las actividades de un usuario
	router.GET("/users/:id", controllers.VerifyToken, controllers.GetUserByID)                  //trae un usuario por id
	router.POST("/activity", controllers.VerifyAdminToken, controllers.CreateActivity)          //crea una actividad
	router.DELETE("/activity/:id", controllers.VerifyAdminToken, controllers.DeleteActivity)
	router.PUT("/activity/:id", controllers.VerifyAdminToken, controllers.ModifyActivity)
	router.GET("/users/admin", controllers.VerifyAdminToken)
}
