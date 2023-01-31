package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Nacionality string        `json:"nacionality" bson:"nacionality"`
	Age         int           `json:"age" bson:"age"`
}

// creamos una función "NewController", esta se ejecutará en el main
// esta tendrá como parametro una variable tipo mgo.Client el cual será producto de la función GetSession()
type UserController struct {
	mongoClient *mongo.Client
}

// La función NewUserController crea un nuevo controlador de usuario y devuelve un puntero a él.
// Acepta un puntero a una sesión de MongoDB (mgo.Session) como argumento y almacena esa sesión en el controlador
// de usuario devuelto. El controlador de usuario
// es una estructura que probablemente tiene métodos para realizar operaciones en la base
// de datos de MongoDB utilizando esa sesión específica
func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}
func (ctrl UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// investigar la forma de leer un json en el r request
	// crear un struct de forma que puedas pasar el json al struct usando marshall
	// investigar como usar el mongo client para crear un registro en una bd de mongo
	var NewUser User
	err := json.NewDecoder(r.Body).Decode(&NewUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	collection := ctrl.mongoClient.Database("mongo-golang").Collection("users")
	_, err = collection.InsertOne(context.TODO(), NewUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Doy el aviso de un usuario creado
	w.WriteHeader(http.StatusCreated)
}
func (ctrl UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//En primera, tenemos que elegir la ID del usuario a eliminar
	id := p.ByName("id")
	//verifica que la id ingresada sea un ID en base hexadecimal valida, si no lo es , retorna "false"
	if !bson.IsObjectIdHex(id) {
		http.Error(w, "ID no valida", http.StatusBadRequest)
		return
	}
	//
	oid := bson.ObjectIdHex(id)
	//Collection es el conjunto de documento json guardados en nuestra base de datos, collection es un variable puntero de tipo *mongo.Collection,
	// nosotros estamos creando en cierta manera, un acceso...
	collection := ctrl.mongoClient.Database("mongo-golang").Collection("users")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func (ctrl UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	//verifico que la id sea un numero u objeto en sistema hexadecimal
	if !bson.IsObjectIdHex(id) {
		http.Error(w, "ID no valida", http.StatusBadRequest)
		return
	}
	//guardo la ID en el formato bson.ObjectId
	oid := bson.ObjectIdHex(id)
	collection := ctrl.mongoClient.Database("mongo-golang").Collection("users")
	var thatuser User
	//Usamos la función encode paara buscar en la colección el documento con la ID solicitada y transformarla de json a la estructura "User"
	err := collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&thatuser)
	if err != nil {
		// los http.Status... son constantes que indican avisos, en este caso , si el usuario con la ID indicada no es encontrada
		//lanza el famoso error 404
		http.Error(w, err.Error(), 404)
		return
	}
	//caso contrario, mostramos el aviso que la función trabajó adecuadamente
	// los parámetros colocados son estrictamente
	userJSON, err := json.Marshal(thatuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "aplication/json")
	//Programa Indica que la función fue exitosa
	w.WriteHeader(http.StatusOK)
	//muestra El usuario que solicité!
	w.Write(userJSON)
}

func (ctrl UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		http.Error(w, "ID no valida", http.StatusBadRequest)
		return
	}
	oid := bson.ObjectIdHex(id)
	var upuser User
	//decode sirve para transformar ese json en la estructura que deseo almacenarla
	err := json.NewDecoder(r.Body).Decode(&upuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	upuser.ID = oid
	collection := ctrl.mongoClient.Database("mongo-golang").Collection("users")
	_, err = collection.ReplaceOne(context.TODO(), bson.M{"_id": oid}, upuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
