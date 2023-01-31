package mongocent

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient() *mongo.Client {
	//Especifico la versión del API de mi servidor mongoDB
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	//Configura las opciones del cliente con la URI de conexión a la base de datos y las opciones del servidor.
	clienOptions := options.Client().ApplyURI("mongodb+srv://Boluarte616:cacanadedog@elclust3r.kk21bg1.mongodb.net/test").
		SetServerAPIOptions(serverAPIOptions)
	//Acá creo un contexto, un contexto es un objeto que se utiliza para compartir información como tiempo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	/*El propósito de llamar a cancel aquí es cancelar el contexto creado
	con context.WithTimeout. Al cancelar un contexto, se envía una señal a las
	operaciones que se están ejecutando en ese contexto, indicándoles que
	deben detenerse inmediatamente. De esta manera, se garantiza que se
	liberan los recursos utilizados por esas operaciones y se evitan posibles
	errores o leaks de memoria.*/
	defer cancel()

	client, err := mongo.Connect(ctx, clienOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
