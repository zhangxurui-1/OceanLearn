package main

import (
	"github.com/gin-gonic/gin"
	"oceanlearn/common"
)

func main() {
	common.InitDB()
	r := gin.Default()
	r = CollectRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}
