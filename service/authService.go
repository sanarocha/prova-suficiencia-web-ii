package service

import (
	"context"
	"time"

	"github.com/sanarocha/prova-suficiencia-web-ii/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	Collection *mongo.Collection
}

// buscar user no banco de dados a partir do username
func (s *AuthService) FindUser(username string) (*models.User, error) {
	var user models.User
	// instancia contexto com timeout de 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// liberar recursos
	defer cancel()

	// consulta no DB, FindOne retorna único documento da collection com username correspondente e decodifica ele para struct
	err := s.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// cadastra novo usuário
func (s *AuthService) RegisterUser(user *models.User) (*mongo.InsertOneResult, error) {
	// contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insere documento user na collection de usuários do BD
	result, err := s.Collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}
