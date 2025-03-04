package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oceanlearn/model"
	"oceanlearn/repository"
	"oceanlearn/response"
	"oceanlearn/vo"
	"strconv"
)

type ICategoryController interface {
	RestController
}
type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repo := repository.NewCategoryRepository()
	repo.DB.AutoMigrate(&model.Category{})

	return CategoryController{Repository: repo}
}

func (c CategoryController) Create(ctx *gin.Context) {
	// view object
	requestCategory := vo.CreateCategoryRequest{}
	// 数据绑定异常，返回500
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，请输入分类名")
		return
	}

	category, err := c.Repository.Create(requestCategory.Name)
	if err != nil || category == nil {
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"category": *category}, "")
}

func (c CategoryController) Update(ctx *gin.Context) {
	var (
		category   *model.Category
		err        error
		categoryId int
	)
	requestCategory := vo.CreateCategoryRequest{}
	// 数据绑定异常，返回500
	if err = ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，请输入分类名")
		return
	}
	// 从path获取分类id
	categoryId, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "无效的分类ID")
		return
	}

	// 查找
	category, err = c.Repository.SelectById(categoryId)
	if err != nil || category == nil {
		response.Response(ctx, http.StatusNotFound, 404, nil, "分类未找到")
		return
	}

	// 更新
	category, err = c.Repository.Update(category, requestCategory.Name)
	if err != nil || category == nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"category": *category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	var (
		category   *model.Category
		err        error
		categoryId int
	)
	// 从path获取分类id
	categoryId, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "无效的分类ID")
		return
	}

	// 查找
	category, err = c.Repository.SelectById(categoryId)
	if err != nil || category == nil {
		response.Response(ctx, http.StatusNotFound, 404, nil, "分类未找到")
		return
	}

	response.Success(ctx, gin.H{"category": *category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	// 从path获取分类id
	categoryId, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		response.Fail(ctx, nil, "无效的分类ID")
		return
	}

	// 删除
	if err = c.Repository.DeleteById(categoryId); err != nil {
		response.Fail(ctx, nil, "删除失败，请重试")
		return
	}

	response.Success(ctx, nil, "删除成功")
}
