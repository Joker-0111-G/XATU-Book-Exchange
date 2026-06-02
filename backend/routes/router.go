package routes

import (
	"xatu-book-exchange/handler"
	"xatu-book-exchange/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 静态文件（上传的图片）
	r.Static("/uploads", "./uploads")

	// 提供前端入口（后端工作目录为 backend/，前端在 ../frontend/）
	r.GET("/", func(c *gin.Context) {
		c.File("../frontend/index.html")
	})

	// 初始化 handlers
	userHandler := handler.NewUserHandler()
	bookHandler := handler.NewBookHandler()
	categoryHandler := handler.NewCategoryHandler()
	orderHandler := handler.NewOrderHandler()
	favoriteHandler := handler.NewFavoriteHandler()
	messageHandler := handler.NewMessageHandler()
	uploadHandler := handler.NewUploadHandler()
	adminHandler := handler.NewAdminHandler()
	bannerHandler := handler.NewBannerHandler()

	// === 公开接口（无需登录，带限流） ===
	public := r.Group("/api/v1/public")
	public.Use(middleware.RateLimiter(middleware.RateLimitConfig{
		MaxRequests: 60,
		Window:      1 * 60,
	}))
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
		public.POST("/refresh-token", userHandler.RefreshToken)
		public.GET("/books", bookHandler.List)
		public.GET("/books/:id", bookHandler.Get)
		public.GET("/books/search", bookHandler.Search)
		public.GET("/categories", categoryHandler.List)
		public.GET("/banners", bannerHandler.List)
			public.GET("/statistics", bannerHandler.Statistics)
	}

	// === 需要登录的接口 ===
	auth := r.Group("/api/v1")
	auth.Use(middleware.AuthRequired())
	{
		// 用户相关
		auth.POST("/user/logout", userHandler.Logout)
		auth.GET("/user/profile", userHandler.Profile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)
		auth.PUT("/user/password", userHandler.ChangePassword)
		auth.GET("/user/books", bookHandler.UserBooks)
		auth.GET("/user/orders", orderHandler.UserOrders)
		auth.GET("/user/sales", orderHandler.UserSales)

		// 图书管理
		auth.POST("/books", bookHandler.Create)
		auth.PUT("/books/:id", bookHandler.Update)
		auth.PUT("/books/:id/status", bookHandler.UpdateStatus)
		auth.DELETE("/books/:id", bookHandler.Delete)

		// 收藏
		auth.GET("/favorites", favoriteHandler.List)
		auth.POST("/favorites", favoriteHandler.Add)
		auth.DELETE("/favorites/:bookId", favoriteHandler.Remove)
		auth.GET("/favorites/check/:bookId", favoriteHandler.Check)

		// 订单
		auth.POST("/orders", orderHandler.Create)
		auth.GET("/orders", orderHandler.List)
		auth.GET("/orders/sales", orderHandler.SalesList)
		auth.GET("/orders/:id", orderHandler.Get)
		auth.PUT("/orders/:id/confirm", orderHandler.Confirm)
		auth.PUT("/orders/:id/complete", orderHandler.Complete)
		auth.PUT("/orders/:id/cancel", orderHandler.Cancel)

		// 消息
		auth.GET("/messages/conversations", messageHandler.Conversations)
		auth.GET("/messages/conversations/:userId", messageHandler.History)
		auth.POST("/messages", messageHandler.Send)
		auth.PUT("/messages/read", messageHandler.MarkRead)
		auth.GET("/messages/unread-count", messageHandler.UnreadCount)

		// 上传
		auth.POST("/upload/image", uploadHandler.UploadImage)
		auth.DELETE("/upload/image", uploadHandler.DeleteImage)
	}

	// === 管理员接口 ===
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/users", adminHandler.Users)
		admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		admin.GET("/books", adminHandler.Books)
		admin.PUT("/books/:id/status", adminHandler.UpdateBookStatus)
		admin.GET("/orders", adminHandler.Orders)
		admin.GET("/categories", adminHandler.Categories)
		admin.POST("/categories", adminHandler.CreateCategory)
		admin.PUT("/categories/:id", adminHandler.UpdateCategory)
		admin.DELETE("/categories/:id", adminHandler.DeleteCategory)
		admin.POST("/banners", adminHandler.CreateBanner)
		admin.PUT("/banners/:id", adminHandler.UpdateBanner)
		admin.DELETE("/banners/:id", adminHandler.DeleteBanner)
		admin.GET("/statistics", adminHandler.Statistics)
	}

	return r
}
