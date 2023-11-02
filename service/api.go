package service

import (
	"bookmanager-server/global"
	"bookmanager-server/model"
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetApis(currentPage, pageSize int, keyWord string) (gin.H, error) {
	var apis []model.Api
	var count int64
	opts := options.Find().SetSkip(int64((currentPage - 1) * pageSize)).SetLimit(int64(pageSize))
	filter := bson.M{}

	if keyWord != "" {
		regexString := fmt.Sprintf(".*%v.*", keyWord)
		filter = bson.M{
			"$or": []bson.M{
				{
					"name": primitive.Regex{Pattern: regexString},
				},
				{
					"url": primitive.Regex{Pattern: regexString},
				},
				{
					"method": primitive.Regex{Pattern: regexString},
				},
				{
					"desc": primitive.Regex{Pattern: regexString},
				},
			},
		}
	}

	cursor, err := global.ApiColl.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &apis); err != nil {
		fmt.Println(err)
		return nil, err
	}

	count, err = global.ApiColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return gin.H{
		"apis":  apis,
		"count": count,
	}, nil
}

func PostApis(apis []model.Api) (gin.H, error) {
	var existApis []model.Api
	for _, api := range apis {
		filter := bson.M{
			"url":    api.Url,
			"method": api.Method,
		}
		if err := global.ApiColl.FindOne(context.TODO(), filter).Decode(&model.Api{}); errors.Is(err, mongo.ErrNoDocuments) {
			_, err = global.ApiColl.InsertOne(context.TODO(), api)
			if err != nil {
				return nil, err
			}
		} else {
			existApis = append(existApis, api)
		}
	}

	if len(existApis) == 0 {
		return nil, nil
	} else {
		return gin.H{
			"existApis": existApis,
		}, errors.New("existApis")
	}
}

func UpdateApi(api *model.Api) (gin.H, error) {
	// url与method不能相同
	if err := global.ApiColl.FindOne(context.TODO(), bson.M{
		"url":    api.Url,
		"method": api.Method,
	}).Decode(&model.Api{}); !errors.Is(err, mongo.ErrNoDocuments) {
		return gin.H{
			"error": "存在相同的api",
		}, errors.New("存在相同的api")
	}

	filter := bson.M{
		"_id": api.Id,
	}
	update := bson.M{
		"$set": api,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(1)
	var DBApi model.Api

	err := global.ApiColl.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&DBApi)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"api": DBApi,
	}, nil
}

func DeleteApi(Id primitive.ObjectID) (gin.H, error) {
	_, err := global.ApiColl.DeleteOne(context.TODO(), bson.M{"_id": Id})
	if err != nil {
		return nil, err
	}
	if _, err := global.RoleColl.UpdateMany(context.TODO(), bson.M{}, bson.M{
		"$pull": bson.M{
			"apis": Id,
		},
	}); err != nil {
		return gin.H{
			"error": "更新失败",
		}, errors.New("更新失败")
	}
	return nil, nil
}
