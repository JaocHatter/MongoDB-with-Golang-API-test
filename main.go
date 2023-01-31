package main

import (
	"net/http"

	mongocent "github.com/JaocHatter/mongo-golang/MDBClient"
	"github.com/JaocHatter/mongo-golang/controllers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	//Creamos un router
	r := httprouter.New()
	//la funcion NewUserController tomará como parámetro otra función llamada GetClient() del paquete mongocent,
	usercont := controllers.NewUserController(mongocent.GetClient())
	r.POST("/new", usercont.CreateUser)
	r.GET("/create", usercont.GetUser)
	r.DELETE("/delete/:id", usercont.DeleteUser)
	r.PUT("/update/:id", usercont.UpdateUser)
	http.ListenAndServe("localhost:666", r)
}
