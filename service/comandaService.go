package service

import (
	"context"
	"errors"
	"time"

	"github.com/sanarocha/prova-suficiencia-web-ii/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComandaService struct {
	Collection *mongo.Collection
}

func (s *ComandaService) GetComandaByID(id primitive.ObjectID) (*models.ComandaGetIDResponse, error) {
	// variável que vai armazenar dados da comanda encontrada
	var comanda models.ComandaGetIDResponse
	// contexto com timeout de 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// consulta no BD onde ID seja igual ao ID fornecido, se encontrado é decodificado e armazenado
	err := s.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comanda)

	if err != nil {
		return nil, err
	}

	// retorna comanda e nenhum erros
	return &comanda, nil
}

func (s *ComandaService) CreateComanda(comanda *models.ComandaInput) (*mongo.InsertOneResult, error) {
	// context com timeout de 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insere comanda no BD
	// retorno do InsertOne é o ID do objeto criado
	return s.Collection.InsertOne(ctx, comanda)
}

func (s *ComandaService) UpdateComanda(id primitive.ObjectID, updateData models.ComandaInput) (*models.ComandaGetIDResponse, error) {
	// context com timeout de 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mapa BSON vazio, preenchido com os campos que forem ser atualizados
	update := bson.M{}

	// verifica se campo não é vazio, senão for é adicionado ao mapa update
	if updateData.IDUsuario != nil {
		update["idUsuario"] = *updateData.IDUsuario
	}
	if updateData.NomeUsuario != nil && *updateData.NomeUsuario != "" {
		update["nomeUsuario"] = *updateData.NomeUsuario
	}
	if updateData.TelefoneUsuario != nil && *updateData.TelefoneUsuario != "" {
		update["telefoneUsuario"] = *updateData.TelefoneUsuario
	}

	// valida se há produtos e se os campos são válidos
	if updateData.Produtos != nil && len(*updateData.Produtos) > 0 {
		validProducts := []models.ProdutoInput{}
		for _, produto := range *updateData.Produtos {
			if produto.Nome != nil && *produto.Nome != "" && produto.Preco != nil && *produto.Preco != 0 {
				// armazena em uma lista os produtos válidos
				validProducts = append(validProducts, produto)
			}
		}
		// associa os produtos válidos ao mapa update
		if len(validProducts) > 0 {
			update["produtos"] = validProducts
		}
	}

	// se mapa estiver vazio, retorna erro
	if len(update) == 0 {
		return nil, errors.New("Nenhum campo para atualizar")
	}

	// executa operação de atualização
	_, err := s.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	// retorna comanda atualizada através do método GetComandaByID
	return s.GetComandaByID(id)
}

func (s *ComandaService) DeleteComanda(id primitive.ObjectID) error {
	// contexto com timeout de 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// operação de remoção no banco
	result, err := s.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	// verifica se alguma comanda foi deletada
	if result.DeletedCount == 0 {
		// retorna erro se nenhuma comanda foi removida
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *ComandaService) GetAllComandas() ([]models.ComandaGetResponse, error) {
	// slice vazio de comandas
	var comandas []models.ComandaGetResponse
	// contexto com timeout 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// operação de busca com filtro vazio, método Find retorna um cursor que é usado para iteração
	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// garante que o cursor será fechado
	defer cursor.Close(ctx)

	// iteração sobre os dados retornados na consulta
	for cursor.Next(ctx) {
		var comanda models.ComandaGetResponse
		// decodificar comandas para a model e armazenar na slice de retorno
		if err = cursor.Decode(&comanda); err != nil {
			return nil, err
		}
		comandas = append(comandas, comanda)
	}

	return comandas, nil
}
