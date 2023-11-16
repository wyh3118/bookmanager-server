package middleware

import (
	"bookmanager-server/global"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

// UserClaims 在token内额外记录用户的id
type userClaims struct {
	Id primitive.ObjectID
	jwt.RegisteredClaims
}

// tokenExpireDuration token的过期时间
const tokenExpireDuration = time.Hour * 24

// userSecret 密钥使用UUID
var userSecret = []byte("eae887e5-f1ac-418e-b54a-14ef5cf6f480")

func CreateToken(Id primitive.ObjectID) (tokenString string, err error) {
	// 创建一个我们自己的声明
	claims := userClaims{
		Id, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)),
			Issuer:    "wyh", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err = token.SignedString(userSecret)
	// 将tokenString存入redis
	global.Redis.Set(context.TODO(), "user:"+Id.Hex(), tokenString, 0)
	return tokenString, err
}

func parseToken(tokenString string) (*userClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return userSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*userClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "没有token",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "token格式错误",
			})
			c.Abort()
			return
		}
		tokenString := parts[1]
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		uc, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "无效的Token",
			})
			c.Abort()
			return
		}
		//判断tokenString是否与redis中的相同以保证只能有一个设备登录
		redisToken := global.Redis.Get(context.TODO(), "user:"+uc.Id.Hex()).Val()
		if redisToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "账号已在其它地方登录，请重新登录",
			})
			c.Abort()
			return
		}
		// 将当前请求的Id信息保存到请求的上下文c上
		c.Set("Id", uc.Id)
		c.Next() // 后续的处理函数可以用过c.Get("Id")来获取当前请求的用户Id
	}
}
