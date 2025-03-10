package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"oceanlearn/common"
	"oceanlearn/model"
	"strings"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestMap := model.User{}
		if err := ctx.Bind(&requestMap); err != nil {
			return
		}
		// 获取authorization header
		tokenstring := ctx.GetHeader("Authorization")

		if tokenstring == "" || !strings.HasPrefix(tokenstring, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			ctx.Abort()
			return
		}

		// 查 redis 缓存
		telephone := requestMap.Telephone
		redisClient := common.GetRedisClient()
		tokenCache, err := redisClient.Get(context.Background(), telephone+":token").Result()
		// redis 中没有找到或 token 不一致
		if err != nil || tokenCache != tokenstring {
			// 剔除 token 头 "Bearer "
			tokenstring = tokenstring[7:]

			token, claims, err2 := common.ParseToken(tokenstring)
			if err2 != nil || !token.Valid {
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

			// 存入Redis
			redisClient.Set(context.Background(), telephone+":token", "Bearer "+tokenstring, time.Hour)
			redisClient.Set(context.Background(), telephone+":name", user.Name, time.Hour)

			ctx.Set("telephone", telephone)
		}

		ctx.Next()
	}
}
