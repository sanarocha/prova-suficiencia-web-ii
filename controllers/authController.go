package controllers

import (
	"net/http"

	"github.com/sanarocha/prova-suficiencia-web-ii/middlewares"
	"github.com/sanarocha/prova-suficiencia-web-ii/models"
	"github.com/sanarocha/prova-suficiencia-web-ii/service"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var authService service.AuthService

// configurando o authService para se conectar com a collection usuarios do BD
func InitializeAuth(client *mongo.Client) {
	authService = service.AuthService{
		Collection: client.Database("bdProvaSuficiencia").Collection("usuarios"),
	}
}

// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.User true "Token de login"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {object} map[string]string "erro"
// @Failure      401 {object} map[string]string "erro"
// @Router       /login [post]
func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// vincular dados recebidos do json a struct credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Os dados inseridos são inválidos"})
		return
	}

	// verificação se o usuário existe através de função auxiliar da camada service
	user, err := authService.FindUser(credentials.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Nome de usuário ou senha incorretos"})
		return
	}

	// verificação se a senha é válida através do método bcrypt
	// compara a senha fornecida com o hash armazenado no BD
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Nome de usuário ou senha incorretos"})
		return
	}

	// gera o token JWT através do middleware passando o user ID como parte da claims
	token, err := middlewares.CreateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao gerar token"})
		return
	}

	// retorna o token se sucesso
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary      Cadastro
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "Dados para cadastro"
// @Success      201 {object} map[string]string "mensagem"
// @Failure      400 {object} map[string]string "erro"
// @Failure      409 {object} map[string]string "erro"
// @Failure      500 {object} map[string]string "erro"
// @Router       /register [post]
func Register(c *gin.Context) {
	var user models.User
	// vincular dados recebidos com a model de usuário
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Os dados inseridos são inválidos"})
		return
	}

	// verificação se o usuário existe através de função auxiliar da camada service
	// se já existir retorna erro
	if existingUser, _ := authService.FindUser(user.Username); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "Nome de usuário já existe"})
		return
	}

	// criptografa a senha fornecida com bcrypt com custo padrão
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao cadastrar senha"})
		return
	}

	// armazena a senha a model do usuário
	user.Password = string(hashedPassword)

	// cadastrar o usuário através da camada service
	_, err = authService.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Não foi possível completar o cadastro"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário criado com sucesso"})
}
