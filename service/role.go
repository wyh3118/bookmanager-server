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

func GetRoles(currentPage, pageSize int, keyWord string) (gin.H, error) {
	var roles []model.Role
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
					"desc": primitive.Regex{Pattern: regexString},
				},
			},
		}
	}

	cursor, err := global.RoleColl.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &roles); err != nil {
		fmt.Println(err)
		return nil, err
	}

	count, err = global.RoleColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return gin.H{
		"roles": roles,
		"count": count,
	}, nil
}

func PostRoles(roles []model.Role) (gin.H, error) {
	var existRoles []model.Role

	for _, role := range roles {
		err := global.RoleColl.FindOne(context.TODO(), bson.M{"name": role.Name}).Decode(&model.Role{})
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err = global.RoleColl.InsertOne(context.TODO(), role)
			if err != nil {
				return gin.H{
					"error": "插入失败",
				}, err
			}
		} else {
			existRoles = append(existRoles, role)
		}
	}

	if len(existRoles) == 0 {
		return nil, nil
	} else {
		return gin.H{
			"existRoles": existRoles,
		}, errors.New("存在同名身份")
	}
}

func UpdateRole(role *model.Role) (gin.H, error) {
	err := global.RoleColl.FindOne(context.TODO(), bson.M{"_id": role.Id}).Decode(&model.Role{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return gin.H{
			"error": "没有此角色",
		}, err
	}

	filer := bson.M{
		"_id": role.Id,
	}
	update := bson.M{
		"$set": role,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(1)
	DBRole := model.Role{}
	err = global.RoleColl.FindOneAndUpdate(context.TODO(), filer, update, opts).Decode(&DBRole)
	if err != nil {
		return gin.H{
			"error": "更新失败",
		}, err
	}
	return gin.H{
		"role": DBRole,
	}, nil
}

func DeleteRole(roleId primitive.ObjectID) (gin.H, error) {
	err := global.RoleColl.FindOneAndDelete(context.TODO(), bson.M{"_id": roleId}).Decode(&model.Role{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return gin.H{
			"error": "不存在此角色",
		}, errors.New("不存在此角色")
	}
	return nil, nil
}
