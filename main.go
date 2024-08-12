package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sanarocha/prova-suficiencia-web-ii/controllers"
	_ "github.com/sanarocha/prova-suficiencia-web-ii/docs"
	"github.com/sanarocha/prova-suficiencia-web-ii/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// @title           Swagger API Documentation
// @version         1.0
// @description     Documentação da API para a prova de suficiência da disciplina Programação Web II
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// configuração da conexão com o MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// conectando ao MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// verificando a conexão
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Conectado ao banco de dados")

	// inicializando os serviços de comanda e auth
	controllers.Initialize(client)
	controllers.InitializeAuth(client)

	// configuração do Gin
	r := gin.Default()

	// configuração das rotas e inicialização do servidor
	r = routes.SetupRouter()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// iniciando o servidor na porta 8080
	r.Run(":8080")
}
