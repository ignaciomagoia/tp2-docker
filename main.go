package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

var userCollection *mongo.Collection

func main() {
	// Leer variable de entorno
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI no está definida en las variables de entorno")
	}

	// Conexión a MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Usar base de datos y colección
	db := client.Database("hotelapp")
	userCollection = db.Collection("users")

	// Iniciar Gin
	r := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// Registro de usuario
	r.POST("/register", registerUser)

	// Login de usuario
	r.POST("/login", loginUser)

	// Endpoints de testing
	r.GET("/users", listUsers)
	r.DELETE("/users", clearUsers)

	// Iniciar servidor
	r.Run(":8080")
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Verificar si ya existe el usuario
	var existing User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Usuario ya existe"})
		return
	}

	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado con éxito"})
}

func loginUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var found User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&found)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar password
	if found.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password incorrecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso"})
}

func listUsers(c *gin.Context) {
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	defer cursor.Close(context.TODO())

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar usuarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func clearUsers(c *gin.Context) {
	_, err := userCollection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al limpiar usuarios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todos los usuarios han sido eliminados"})
}
