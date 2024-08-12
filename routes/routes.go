package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sanarocha/prova-suficiencia-web-ii/controllers"
	"github.com/sanarocha/prova-suficiencia-web-ii/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// rotas de autenticação
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	// rotas de comandas
	comandaRoutes := r.Group("/comandas")
	{
		// rotas GET - sem autenticação
		comandaRoutes.GET("/", controllers.GetAllComandas)
		comandaRoutes.GET("/:id", controllers.GetComandaByID)

		// rotas POST, PUT, DELETE - com autenticação
		comandaRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateComanda)
		comandaRoutes.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateComanda)
		comandaRoutes.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteComanda)
	}

	return r
}
