package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"oceanlearn/common"
	"oceanlearn/dto"
	"oceanlearn/model"
	"oceanlearn/response"
	"oceanlearn/util"
	"time"
)

func Register(ctx *gin.Context) {
	requestMap := model.User{}
	if err := ctx.Bind(&requestMap); err != nil {
		return
	}

	name := requestMap.Name
	telephone := requestMap.Telephone
	password := requestMap.Password

	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不得小于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandString(10)
	}

	if isTelephoneExist(telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已注册")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密异常")
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	DB := common.GetDB()
	DB.Create(&newUser)

	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v\n", err)
		return
	}
	// 存入 redis
	redisClient := common.GetRedisClient()
	redisClient.Set(context.Background(), telephone+":token", token, time.Hour)
	redisClient.Set(context.Background(), telephone+":name", newUser.Name, time.Hour)

	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	requestMap := model.User{}
	if err := ctx.Bind(&requestMap); err != nil {
		return
	}
	telephone := requestMap.Telephone
	password := requestMap.Password

	var user model.User
	DB := common.GetDB()
	DB.Where("telephone = ?", telephone).First(&user)

	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 服务端发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error: %v\n", err)
		return
	}
	// 存入 redis
	redisClient := common.GetRedisClient()
	redisClient.Set(context.Background(), telephone+":token", token, time.Hour)
	redisClient.Set(context.Background(), telephone+":name", user.Name, time.Hour)

	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	requestMap := model.User{}
	if err := ctx.Bind(&requestMap); err != nil {
		return
	}

	// 从 redis 获取用户名；认证中间件已经存过相应的内容
	redisClient := common.GetRedisClient()
	telephone, _ := ctx.Get("telephone")
	username, err := redisClient.Get(context.Background(), telephone.(string)+":name").Result()
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "服务异常")
		return
	}

	response.Success(ctx, gin.H{
		"user": dto.UserDto{
			Name:      username,
			Telephone: telephone.(string),
		},
	}, "")
}

func isTelephoneExist(telephone string) bool {
	var user model.User
	DB := common.GetDB()
	DB.Where("telephone = ?", telephone).First(&user)

	if user.ID != 0 {
		return true
	}
	return false
}
