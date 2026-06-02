package repository

import (
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
)

type MessageRepo struct{}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{}
}

func (r *MessageRepo) Create(msg *model.Message) error {
	return database.DB.Create(msg).Error
}

func (r *MessageRepo) FindConversations(userID uint) ([]model.Message, error) {
	// 使用 RAW SQL 获取每个会话的最新消息（带参数，防注入）
	subQuery := database.DB.Raw(`
		SELECT MAX(m.id)
		FROM messages m
		WHERE ? IN (m.from_user_id, m.to_user_id)
		GROUP BY CASE WHEN m.from_user_id = ? THEN m.to_user_id ELSE m.from_user_id END
	`, userID, userID)

	var msgs []model.Message
	err := database.DB.Where("id IN (?)", subQuery).
		Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC").
		Find(&msgs).Error
	return msgs, err
}

func (r *MessageRepo) FindHistory(userID, otherUserID uint) ([]model.Message, error) {
	var msgs []model.Message
	err := database.DB.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userID, otherUserID, otherUserID, userID,
	).Order("created_at ASC").Find(&msgs).Error
	return msgs, err
}

func (r *MessageRepo) MarkRead(fromUserID, toUserID uint) error {
	return database.DB.Model(&model.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = 0", fromUserID, toUserID).
		Update("is_read", 1).Error
}

func (r *MessageRepo) UnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Message{}).
		Where("to_user_id = ? AND is_read = 0", userID).
		Count(&count).Error
	return count, err
}
