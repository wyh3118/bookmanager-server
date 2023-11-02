package service

import (
	"bookmanager-server/global"
	"bookmanager-server/middleware"
	"bookmanager-server/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const defaultHead = "https://img1.baidu.com/it/u=3875084872,415645148&fm=253&fmt=auto&app=138&f=JPEG?w=380&h=380"

var simpleUser, _ = primitive.ObjectIDFromHex("6541fd9a1950919844415948")

func CreateUser(user *model.User) (gin.H, error) {
	user.Head = defaultHead

	// 判断邮箱是否重复
	err := global.UserColl.FindOne(context.TODO(), bson.M{
		"email": user.Email,
	}).Decode(&model.User{})

	// 查找后没有error说明存在相同的邮箱
	if err == nil {
		return gin.H{
			"error": "邮箱已存在",
		}, err
	}

	// 对用户密码加密
	newPassWord, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(newPassWord)

	// 身份默认为普通用户
	user.Roles = append(user.Roles, simpleUser)

	// 存入数据库
	insertResult, err := global.UserColl.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	// 生成token
	Id := insertResult.InsertedID.(primitive.ObjectID)
	token, err := middleware.CreateToken(Id)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"user":  user,
		"token": token,
	}, nil
}

func GetUser(user *model.User) (gin.H, error) {
	DBUser := &model.User{}

	// 根据邮箱查找用户
	err := global.UserColl.FindOne(context.TODO(), bson.M{
		"email": user.Email,
	}).Decode(DBUser)
	if err != nil {
		return gin.H{
			"error": "邮箱不存在",
		}, err
	}

	// 密码校验
	if err = bcrypt.CompareHashAndPassword([]byte(DBUser.Password), []byte(user.Password)); err != nil {
		return gin.H{
			"error": "密码错误",
		}, err
	}

	// 生成token
	token, err := middleware.CreateToken(DBUser.Id)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"user":  DBUser,
		"token": token,
	}, nil
}

func UpdateUser(user *model.User) (gin.H, error) {
	updateData := make(map[string]string)

	if user.Username != "" {
		updateData["username"] = user.Username
	}
	if user.Head != "" {
		updateData["head"] = user.Head
	}
	if user.Password != "" {
		byteData, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		hashedPassword := string(byteData)
		updateData["password"] = hashedPassword
	}

	update := bson.M{
		"$set": updateData,
	}

	_, err := global.UserColl.UpdateByID(context.TODO(), user.Id, update)
	if err != nil {
		return nil, err
	}

	updatedUser := &model.User{}
	err = global.UserColl.FindOne(context.TODO(), bson.M{"_id": user.Id}).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"user": updatedUser,
	}, nil
}

func UpdateUserByAdmin(user *model.User) (gin.H, error) {
	updateData := make(map[string]interface{})
	if user.Username != "" {
		updateData["username"] = user.Username
	}
	if user.Head != "" {
		updateData["head"] = user.Head
	}
	if user.Password != "" {
		byteData, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		hashedPassword := string(byteData)
		updateData["password"] = hashedPassword
	}
	if len(user.Roles) != 0 {
		updateData["roles"] = user.Roles
	}

	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": updateData}
	opts := options.FindOneAndUpdate().SetReturnDocument(1)
	var DBUser model.User

	err := global.UserColl.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&DBUser)
	if err != nil {
		return gin.H{
			"error": "更新失败",
		}, errors.New("更新失败")
	}

	return gin.H{
		"user": DBUser,
	}, nil
}
