package repository

import (
	"errors"
	"gorm.io/gorm"
	"oceanlearn/common"
	"oceanlearn/model"
	"oceanlearn/vo"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository() PostRepository {
	db := common.GetDB()
	return PostRepository{DB: db}
}

func (p *PostRepository) Create(user *model.User, createPostReq *vo.CreatePostRequest) (*model.Post, error) {
	// 创建新 post
	newPost := model.Post{
		UserId:     user.ID,
		CategoryId: createPostReq.CategoryId,
		Title:      createPostReq.Title,
		Content:    createPostReq.Content,
		HeadImgUrl: createPostReq.HeadImgUrl,
	}

	if err := p.DB.Create(&newPost).Error; err != nil {
		return nil, err
	}
	return &newPost, nil
}

func (p *PostRepository) SelectById(id string) (*model.Post, error) {
	// 查询
	post := model.Post{}
	if err := p.DB.Preload("Category").Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepository) Update(post *model.Post, updatePostReq *vo.CreatePostRequest) (*model.Post, error) {
	// 更新
	if err := p.DB.Model(post).Updates(*updatePostReq).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostRepository) DeleteById(id string) error {
	res := p.DB.Where("id = ?", id).Delete(&model.Post{})
	if res.Error != nil {
		// 处理数据库错误
		return res.Error
	} else if res.RowsAffected == 0 {
		// 记录不存在
		return errors.New("record not found")
	}
	return nil
}

// SelectPagingCreatDesc 分页查找最新文章
func (p *PostRepository) SelectPagingCreatDesc(pageNum int, pageSize int) []model.Post {
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)
	return posts
}

// TotalNums 返回记录总数
func (p *PostRepository) TotalNums() int64 {
	var total int64
	p.DB.Model(model.Post{}).Count(&total)
	return total
}
