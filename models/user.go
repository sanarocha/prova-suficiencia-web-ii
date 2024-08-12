package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}
