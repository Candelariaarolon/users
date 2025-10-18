package main

import (
	//importo modulo propio
	"backend/app" //importo modulo propio
	"backend/db"  //importo modulo propio
	_ "fmt"       //importo libreria externa

	_ "github.com/gin-gonic/gin" //importo un link
)

func main() {

	//variable que me apunta al llamado

	db.StartDbEngine()
	app.StartRoute()

	//el segundo parametro que recibe la funcion Get es la declaracion de una funcion, osea no se ejecutara en ese momento
	//la funcion GetHotel es lo que va a hacer cuando se produzca ese llamado, es una referencia a la funcion, ya que no pasamos parametros

}
