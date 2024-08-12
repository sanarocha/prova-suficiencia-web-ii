package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComandaInput struct {
	ID              *primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	IDUsuario       *int                `json:"idUsuario" bson:"idUsuario"`
	NomeUsuario     *string             `json:"nomeUsuario" bson:"nomeUsuario"`
	TelefoneUsuario *string             `json:"telefoneUsuario" bson:"telefoneUsuario"`
	Produtos        *[]ProdutoInput     `json:"produtos" bson:"produtos"`
}

type ProdutoInput struct {
	ID    *int     ` json:"id" bson:"id"`
	Nome  *string  `json:"nome" bson:"nome"`
	Preco *float64 `json:"preco" bson:"preco"`
}

// --

type ComandaGetResponse struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	IDUsuario       int                `json:"idUsuario" bson:"idUsuario"`
	NomeUsuario     string             `json:"nomeUsuario" bson:"nomeUsuario"`
	TelefoneUsuario string             `json:"telefoneUsuario" bson:"telefoneUsuario"`
}

// --

type ComandaGetIDResponse struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"-"`
	IDUsuario       int                    `json:"idUsuario" bson:"idUsuario"`
	NomeUsuario     string                 `json:"nomeUsuario" bson:"nomeUsuario"`
	TelefoneUsuario string                 `json:"telefoneUsuario" bson:"telefoneUsuario"`
	Produtos        []ProdutoGetIDResponse `json:"produtos" bson:"produtos"`
}

type ProdutoGetIDResponse struct {
	ID    int     `json:"id" bson:"id"`
	Nome  string  `json:"nome" bson:"nome"`
	Preco float64 `json:"preco" bson:"preco"`
}

// --

type ComandaPostResponse struct {
	ID              primitive.ObjectID    `bson:"_id" json:"id"`
	IDUsuario       int                   `json:"idUsuario" bson:"idUsuario"`
	NomeUsuario     string                `json:"nomeUsuario" bson:"nomeUsuario"`
	TelefoneUsuario string                `json:"telefoneUsuario" bson:"telefoneUsuario"`
	Produtos        []ProdutoPostResponse `json:"produtos" bson:"produtos"`
}

type ProdutoPostResponse struct {
	ID    int     `json:"id" bson:"id"`
	Nome  string  `json:"nome" bson:"nome"`
	Preco float64 `json:"preco" bson:"preco"`
}
