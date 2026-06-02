package service

import (
	"xatu-book-exchange/model"
	"xatu-book-exchange/repository"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessageService() *MessageService {
	return &MessageService{
		repo: repository.NewMessageRepo(),
	}
}

type SendMessageReq struct {
	ToUserID uint   `json:"to_user_id" binding:"required"`
	BookID   uint   `json:"book_id"`
	Content  string `json:"content" binding:"required"`
}

func (s *MessageService) Send(fromUserID uint, req *SendMessageReq) (*model.Message, error) {
	msg := &model.Message{
		FromUserID:  fromUserID,
		ToUserID:    req.ToUserID,
		BookID:      req.BookID,
		Content:     req.Content,
		ContentType: 0,
		IsRead:      0,
	}
	if err := s.repo.Create(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageService) Conversations(userID uint) ([]model.Message, error) {
	return s.repo.FindConversations(userID)
}

func (s *MessageService) History(userID, otherUserID uint) ([]model.Message, error) {
	return s.repo.FindHistory(userID, otherUserID)
}

func (s *MessageService) MarkRead(fromUserID, toUserID uint) error {
	return s.repo.MarkRead(fromUserID, toUserID)
}

func (s *MessageService) UnreadCount(userID uint) (int64, error) {
	return s.repo.UnreadCount(userID)
}
