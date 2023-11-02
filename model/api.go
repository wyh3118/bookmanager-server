package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Api struct {
	Id     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Url    string             `json:"url" bson:"url"`
	Method string             `json:"method" bson:"method"`
	Desc   string             `json:"desc" bson:"desc"`
}
