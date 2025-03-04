package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oceanlearn/common"
	"oceanlearn/model"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenstring := ctx.GetHeader("Authorization")

		if tokenstring == "" || !strings.HasPrefix(tokenstring, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			ctx.Abort()
			return
		}

		tokenstring = tokenstring[7:]

		token, claims, err := common.ParseToken(tokenstring)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			ctx.Abort()
			return
		}

		// 根据claim中的userid查表
		DB := common.GetDB()
		var user model.User
		DB.First(&user, claims.UserId)

		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
