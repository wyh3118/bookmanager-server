package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string               `json:"email" bson:"email"`
	Username string               `json:"username" bson:"username"`
	Password string               `json:"password" bson:"password"`
	Head     string               `json:"head" bson:"head"`
	Roles    []primitive.ObjectID `json:"roles,omitempty" bson:"roles,omitempty"`
}
