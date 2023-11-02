package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoClient *mongo.Client
	UserColl    *mongo.Collection
	BookColl    *mongo.Collection
	ApiColl     *mongo.Collection
	RoleColl    *mongo.Collection
)
