package service

import (
	"errors"

	"xatu-book-exchange/model"
	"xatu-book-exchange/repository"
	"xatu-book-exchange/utils"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepo(),
	}
}

type RegisterReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Major    string `json:"major"`
	Wechat   string `json:"wechat"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *model.User `json:"user"`
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Major    string `json:"major"`
	Wechat   string `json:"wechat"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterReq) (*model.User, error) {
	// 检查手机号是否已注册
	existing, err := s.repo.FindByPhone(req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("手机号已注册")
	}

	// 加密密码
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 默认昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = "书友" + req.Phone[len(req.Phone)-4:]
	}

	user := &model.User{
		Phone:        req.Phone,
		PasswordHash: hash,
		Nickname:     nickname,
		Major:        req.Major,
		Wechat:       req.Wechat,
		Email:        req.Email,
		Status:       1,
		IsAdmin:      0,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginReq) (*LoginResp, error) {
	user, err := s.repo.FindByPhone(req.Phone)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if user.Status == 0 {
		return nil, ErrUserDisabled
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredential
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken 刷新 Token（用 refresh_token 换新 token 对）
func (s *UserService) RefreshToken(refreshTokenStr string) (*LoginResp, error) {
	claims, err := utils.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("refresh_token 无效或已过期")
	}

	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status == 0 {
		return nil, ErrUserDisabled
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Profile 获取用户信息
func (s *UserService) Profile(userID uint) (*model.User, error) {
	return s.repo.FindByID(userID)
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileReq) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Major != "" {
		user.Major = req.Major
	}
	if req.Wechat != "" {
		user.Wechat = req.Wechat
	}

	return s.repo.Update(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPwd, newPwd string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(oldPwd, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	hash, err := utils.HashPassword(newPwd)
	if err != nil {
		return err
	}

	user.PasswordHash = hash
	return s.repo.Update(user)
}
