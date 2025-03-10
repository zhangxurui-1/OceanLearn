package controller

import (
	"github.com/gin-gonic/gin"
	"oceanlearn/model"
	"oceanlearn/repository"
	"oceanlearn/response"
	"oceanlearn/vo"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	Repository repository.PostRepository
}

func NewPostController() IPostController {
	repo := repository.NewPostRepository()
	return PostController{Repository: repo}
}

func (p PostController) Create(ctx *gin.Context) {
	var createPostReq vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&createPostReq); err != nil {
		response.Fail(ctx, nil, "数据验证失败")
		return
	}

	// 获取登录用户 user
	// 路由中添加了 AuthMiddleware 中间件，会把 user 添加到 ctx 中
	user, success := ctx.Get("user")
	if !success {
		response.Fail(ctx, nil, "登录信息验证失败")
		return
	}
	// 创建
	usr := user.(model.User)
	if _, err := p.Repository.Create(&usr, &createPostReq); err != nil {
		response.Fail(ctx, nil, "创建失败")
		return
	}
	response.Success(ctx, nil, "创建成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 从 path 获取文章 id
	postId := ctx.Params.ByName("id")

	// 查找
	post, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	response.Success(ctx, gin.H{"post": *post}, "")
}
func (p PostController) Update(ctx *gin.Context) {
	var updatePostReq vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&updatePostReq); err != nil {
		response.Fail(ctx, nil, "数据验证失败")
		return
	}
	// 从 path 获取文章 id
	postId := ctx.Params.ByName("id")

	// 获取登录用户 user
	// 路由中添加了 AuthMiddleware 中间件，会把 user 添加到 ctx 中
	user, success := ctx.Get("user")
	if !success {
		response.Fail(ctx, nil, "登录信息验证失败")
		return
	}

	// 查找
	post, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断用户权限，用户是否为文章作者
	if post.UserId != user.(model.User).ID {
		response.Fail(ctx, nil, "非法操作：不是您的文章")
		return
	}

	// 更新
	updatedPost, err := p.Repository.Update(post, &updatePostReq)
	if err != nil {
		panic(err)
		return
	}

	response.Success(ctx, gin.H{"post": *updatedPost}, "更新成功")
}
func (p PostController) Delete(ctx *gin.Context) {
	// 从 path 获取文章 id
	postId := ctx.Params.ByName("id")

	// 获取登录用户 user
	// 路由中添加了 AuthMiddleware 中间件，会把 user 添加到 ctx 中
	user, success := ctx.Get("user")
	if !success {
		response.Fail(ctx, nil, "登录信息验证失败")
		return
	}
	// 查找
	post, err := p.Repository.SelectById(postId)
	if err != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	// 判断用户权限，用户是否为文章作者
	if post.UserId != user.(model.User).ID {
		response.Fail(ctx, nil, "非法操作：不是您的文章")
		return
	}
	// 删除
	if err = p.Repository.DeleteById(postId); err != nil {
		panic(err)
		return
	}
	response.Success(ctx, nil, "删除成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	posts := p.Repository.SelectPagingCreatDesc(pageNum, pageSize)

	// 获取文章总数
	postsNum := p.Repository.TotalNums()

	response.Success(ctx, gin.H{"posts": posts, "posts_num": postsNum}, "")
}
