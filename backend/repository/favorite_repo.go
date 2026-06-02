package repository

import (
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
)

type FavoriteRepo struct{}

func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{}
}

func (r *FavoriteRepo) Create(fav *model.Favorite) error {
	return database.DB.Create(fav).Error
}

func (r *FavoriteRepo) Delete(userID, bookID uint) error {
	return database.DB.Where("user_id = ? AND book_id = ?", userID, bookID).
		Delete(&model.Favorite{}).Error
}

func (r *FavoriteRepo) FindByUserID(userID uint) ([]model.Favorite, error) {
	var favs []model.Favorite
	err := database.DB.Where("user_id = ?", userID).
		Preload("Book").
		Preload("Book.Category").
		Preload("Book.User").
		Order("created_at DESC").Find(&favs).Error
	return favs, err
}

func (r *FavoriteRepo) Check(userID, bookID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&model.Favorite{}).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Count(&count).Error
	return count > 0, err
}
