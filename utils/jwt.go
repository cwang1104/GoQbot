package utils

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

//JWT = JSON WEB TOKEN 是一个开放标准，用于作为json对象，在各个地方安全的传输信息
//此信息可以被验证和信任

type JWT struct {
	SigningKey []byte //自定义密钥
}

//NewJWT 初始化
func NewJWT() *JWT {
	//这里要注意，SigningKey 这个值，需要自定义
	return &JWT{SigningKey: []byte("gqbot")}
}

//自定义信息结构，根据需求填写
type Claims struct {
	UID      int //用户id
	UserName string
	jwt.StandardClaims
}

//定义错误信息
var (
	TokenInvalid = errors.New("Token 无效")
)

//CreateToken 创建 token
func (j *JWT) CreateToken(uid int, userName string) (string, error) {
	c := Claims{
		UID:      uid,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,         //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*365, //签名过期时间
			Issuer:    "gqbot",                          //签名颁发者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.SigningKey)
}

//ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, TokenInvalid
	}
	//解析到Claims 构造中
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, TokenInvalid
}

//RefreshToken 更新 token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		c.StandardClaims.ExpiresAt = time.Now().Add(365 * 24 * time.Hour).Unix()
		return j.CreateToken(c.UID, c.UserName)
	}
	return "", TokenInvalid
}
