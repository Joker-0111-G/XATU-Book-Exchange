package service

import (
	"xatu-book-exchange/model"
	"xatu-book-exchange/repository"
)

type FavoriteService struct {
	favRepo  *repository.FavoriteRepo
	bookRepo *repository.BookRepo
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		favRepo:  repository.NewFavoriteRepo(),
		bookRepo: repository.NewBookRepo(),
	}
}

func (s *FavoriteService) Add(userID, bookID uint) error {
	// 检查图书是否存在
	_, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return ErrNotFound
	}

	// 检查是否已收藏
	exists, err := s.favRepo.Check(userID, bookID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExist
	}

	fav := &model.Favorite{
		UserID: userID,
		BookID: bookID,
	}
	return s.favRepo.Create(fav)
}

func (s *FavoriteService) Remove(userID, bookID uint) error {
	return s.favRepo.Delete(userID, bookID)
}

func (s *FavoriteService) List(userID uint) ([]model.Favorite, error) {
	return s.favRepo.FindByUserID(userID)
}

func (s *FavoriteService) Check(userID, bookID uint) (bool, error) {
	return s.favRepo.Check(userID, bookID)
}
