package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"oceanlearn/common"
	"oceanlearn/model"
	"oceanlearn/util"
)

func Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "手机号必须为11位", "length": len(telephone)})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "密码不得小于6位"})
		return
	}

	if len(name) == 0 {
		name = util.RandString(10)
	}

	if isTelephoneExist(telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "用户已注册"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "加密异常"})
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	DB := common.GetDB()
	DB.Create(&newUser)

	ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	var user model.User
	DB := common.GetDB()
	DB.Where("telephone = ?", telephone).First(&user)

	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码错误"})
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统异常"})
		log.Printf("token generate error: %v\n", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"token": token,
		},
		"message": "登录成功",
	})
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
