package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"oceanlearn/common"
)

type IPostController interface {
	RestController
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	return PostController{DB: db}
}

func (postController PostController) Create(ctx *gin.Context) {

}
func (postController PostController) Show(ctx *gin.Context)   {}
func (postController PostController) Update(ctx *gin.Context) {}
func (postController PostController) Delete(ctx *gin.Context) {}
