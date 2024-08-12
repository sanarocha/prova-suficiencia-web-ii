package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanarocha/prova-suficiencia-web-ii/models"
	"github.com/sanarocha/prova-suficiencia-web-ii/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var comandaCollection *mongo.Collection
var comandaService service.ComandaService

// configurando o comandaService para se conectar com a collection comandas do BD
func Initialize(client *mongo.Client) {
	comandaService = service.ComandaService{
		Collection: client.Database("bdProvaSuficiencia").Collection("comandas"),
	}
}

// @Summary      Retorna a comanda com o ID inserido
// @Tags         comandas
// @Accept       json
// @Produce      json
// @Param        id path string true "ID da Comanda"
// @Success      200 {object} models.ComandaGetIDResponse
// @Failure      400 {object} map[string]string "erro"
// @Failure      404 {object} map[string]string "erro"
// @Failure      500 {object} map[string]string "erro"
// @Router       /comandas/{id} [get]
func GetComandaByID(c *gin.Context) {
	// obter parâmetro ID da url
	idParam := c.Param("id")
	// converter o valor string para ObjectID
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O ID inserido é inválido"})
		return
	}

	// buscar comanda através de função da camada service
	comanda, err := comandaService.GetComandaByID(objectId)
	if err != nil {
		// documento com id especificado não encontrado
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"erro": "A comanda com o ID informado não foi encontrada"})
		} else {
			// outro tipo de erro
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comanda)
}

// @Summary      Cadastra uma nova comanda
// @Tags         comandas
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        comanda body models.ComandaInput true "Dados da comanda"
// @Success      200 {object} models.ComandaPostResponse "Comanda criada"
// @Failure      400 {object} map[string]string "erro"
// @Failure      500 {object} map[string]string "erro"
// @Router       /comandas [post]
func CreateComanda(c *gin.Context) {
	var comandaInput models.ComandaInput
	// parse dos dados recebidos na requisição com a model
	if err := c.ShouldBindJSON(&comandaInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "A comanda inserida tem formato inválido"})
		return
	}

	// validação dos campos obrigatórios para cadastro de uma comanda
	if (comandaInput.IDUsuario == nil || *comandaInput.IDUsuario == 0) ||
		(comandaInput.NomeUsuario == nil || *comandaInput.NomeUsuario == "") ||
		(comandaInput.TelefoneUsuario == nil || *comandaInput.TelefoneUsuario == "") ||
		(comandaInput.Produtos == nil || len(*comandaInput.Produtos) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Faltam informações"})
		return
	}

	// validação dos campos obrigatórios para cadastro de um produto
	for _, prodInput := range *comandaInput.Produtos {
		if (prodInput.ID == nil || *prodInput.ID == 0) ||
			(prodInput.Nome == nil || *prodInput.Nome == "") ||
			(prodInput.Preco == nil || *prodInput.Preco <= 0) {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "Faltam informações"})
			return
		}
	}

	// criação do objeto da comanda que será cadastrado no BD
	comanda := models.ComandaInput{
		IDUsuario:       comandaInput.IDUsuario,
		NomeUsuario:     comandaInput.NomeUsuario,
		TelefoneUsuario: comandaInput.TelefoneUsuario,
		Produtos:        comandaInput.Produtos,
	}

	// cadastra a comanda através de função da camada service
	result, err := comandaService.CreateComanda(&comanda)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	// prepara a resposta com os dados da comanda criada
	produtosResponse := make([]models.ProdutoPostResponse, len(*comandaInput.Produtos))
	for i, prodInput := range *comandaInput.Produtos {
		produtosResponse[i] = models.ProdutoPostResponse{
			ID:    *prodInput.ID,
			Nome:  *prodInput.Nome,
			Preco: *prodInput.Preco,
		}
	}

	// objeto de resposta da comanda tem o ID (obtido a partir do resultado da inserção)
	comandaResponse := models.ComandaPostResponse{
		ID:              result.InsertedID.(primitive.ObjectID),
		IDUsuario:       *comandaInput.IDUsuario,
		NomeUsuario:     *comandaInput.NomeUsuario,
		TelefoneUsuario: *comandaInput.TelefoneUsuario,
		Produtos:        produtosResponse,
	}

	// retorna model e o status
	c.JSON(http.StatusOK, comandaResponse)
}

// @Summary      Atualiza uma comanda existente
// @Tags         comandas
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID da Comanda"
// @Param        comanda body models.ComandaInput true "Novos dados da comanda"
// @Success      200 {object} models.ComandaGetIDResponse
// @Failure      400 {object} map[string]string "erro"
// @Failure      500 {object} map[string]string "erro"
// @Router       /comandas/{id} [put]
func UpdateComanda(c *gin.Context) {
	// permitir a atualização parcial dos dados de uma comanda com base no ID

	// obter parâmetro ID string
	idParam := c.Param("id")
	// converter para objeto ID
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O ID inserido é inválido"})
		return
	}

	// parse dos dados do JSON para a model
	var updateData models.ComandaInput
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// chama função do service para atualizar comanda
	updatedComanda, err := comandaService.UpdateComanda(objectId, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	// retorna model e o status
	c.JSON(http.StatusOK, updatedComanda)
}

// @Summary      Remove uma comanda existente
// @Tags         comandas
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID da Comanda"
// @Success      200 {object} map[string]string "mensagem"
// @Failure      400 {object} map[string]string "erro"
// @Failure      404 {object} map[string]string "erro"
// @Failure      500 {object} map[string]string "erro"
// @Router       /comandas/{id} [delete]
func DeleteComanda(c *gin.Context) {
	// deve remover uma comanda do banco de dados com base no ID

	// obter parâmetro ID string
	idParam := c.Param("id")

	// converter para objeto ID
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O ID inserido é inválido"})
		return
	}

	// função da chamada service para deletar a partir do objeto ID
	err = comandaService.DeleteComanda(objectId)
	if err != nil {
		// caso documento com ID especificado não for encontrado
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"erro": "Comanda não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro": err.Error()})
		}
		return
	}

	// retorna mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"mensagem": "Comanda removida"})
}

// @Summary      Retorna todas as comandas
// @Tags         comandas
// @Accept       json
// @Produce      json
// @Success      200 {array} models.ComandaGetResponse
// @Failure      500 {object} map[string]string "erro"
// @Router       /comandas [get]
func GetAllComandas(c *gin.Context) {
	// retornar todas as comandas armazenadas a partir de função da camada service
	comandas, err := comandaService.GetAllComandas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	// status de sucesso e comandas
	c.JSON(http.StatusOK, comandas)
}
