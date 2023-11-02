package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	Id         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Apis       []primitive.ObjectID `json:"apis" bson:"apis"`
	FirstPage  string               `json:"firstPage" bson:"firstPage"`
	RoleRoutes string               `json:"roleRoutes" bson:"roleRoutes"`
	Desc       string               `json:"desc" bson:"desc"`
}
