package repository

import (
	"errors"
	"gorm.io/gorm"
	"oceanlearn/common"
	"oceanlearn/model"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	db := common.DB
	return CategoryRepository{db}
}

func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{Name: name}

	if err := c.DB.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (c CategoryRepository) Update(category *model.Category, name string) (*model.Category, error) {
	if err := c.DB.Model(&category).Update("name", name).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (c CategoryRepository) SelectById(id int) (*model.Category, error) {
	category := model.Category{}
	if err := c.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (c CategoryRepository) DeleteById(id int) error {
	if rowsAffected := c.DB.Delete(&model.Category{}, id).RowsAffected; rowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}
