package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Author  []string           `json:"author" bson:"author"`
	Press   string             `json:"press" bson:"press"`
	Cover   string             `json:"cover" bson:"cover"`
	Lend    bool               `json:"lend" bson:"lend"`
	LogDate string             `json:"logDate" bson:"logDate"`
}
