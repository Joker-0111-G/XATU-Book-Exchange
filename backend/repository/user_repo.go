package repository

import (
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// FindByPhone 根据手机号查找用户
func (r *UserRepo) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

// Update 更新用户信息
func (r *UserRepo) Update(user *model.User) error {
	return database.DB.Save(user).Error
}
