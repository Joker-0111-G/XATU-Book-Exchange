package database

import (
	"log"

	"xatu-book-exchange/model"

	"golang.org/x/crypto/bcrypt"
)

// SeedData 初始化种子数据（仅在表为空时写入）
func SeedData() {
	seedAdmin()
	seedCategories()
	seedBanners()
}

func seedAdmin() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("G20050111g"), bcrypt.DefaultCost)
	if err != nil {
		hash = []byte("$2a$10$Ss4qX.dk43Phadw5/hqY3.SUF8gRRbwK6FVFAcIbC1yMEoYIYtqqW")
	}

	admin := &model.User{
		Phone:        "19992468036",
		PasswordHash: string(hash),
		Nickname:     "Joker-0111",
		Status:       1,
		IsAdmin:      1,
	}
	if err := DB.Create(admin).Error; err != nil {
		log.Printf("⚠️  创建管理员失败: %v", err)
		return
	}
	log.Println("✅ 初始管理员已创建（19992468036 / G20050111g）")
}

func seedCategories() {
	var count int64
	DB.Model(&model.Category{}).Count(&count)
	if count > 0 {
		return
	}

	// 扁平的分类列表：[name, parentID or 0]
	type catDef struct {
		Name     string
		ParentID uint
	}
	cats := []catDef{
		// === 一级：专业大类 ===
		{Name: "计算机科学与工程学院"},
		{Name: "机电工程学院"},
		{Name: "电子信息工程学院"},
		{Name: "经济管理学院"},
		{Name: "建筑工程学院"},
		{Name: "外国语学院"},
		{Name: "理学院"},
		{Name: "材料与化工学院"},
		{Name: "艺术与传媒学院"},
		{Name: "马克思主义学院"},
		{Name: "公共课（跨专业）"},
		{Name: "考研专区 ✦"},
		{Name: "其他"},
	}

	// 先创建一级分类，记录 ID 映射
	var parentIDs []uint
	for _, c := range cats {
		cat := &model.Category{Name: c.Name}
		DB.Create(cat)
		parentIDs = append(parentIDs, cat.ID)
	}

	// 二级分类（子分类）
	subCats := []struct {
		name     string
		parentID uint
	}{
		// 计算机科学与工程学院
		{"程序设计", parentIDs[0]},
		{"数据结构", parentIDs[0]},
		{"操作系统", parentIDs[0]},
		{"计算机网络", parentIDs[0]},
		// 机电工程学院
		{"机械制图", parentIDs[1]},
		{"理论力学", parentIDs[1]},
		{"材料力学", parentIDs[1]},
		// 电子信息工程学院
		{"电路分析", parentIDs[2]},
		{"模拟电子", parentIDs[2]},
		{"数字电子", parentIDs[2]},
		// 经济管理学院
		{"管理学", parentIDs[3]},
		{"微观经济学", parentIDs[3]},
		{"会计学", parentIDs[3]},
		// 建筑工程学院
		{"土木工程材料", parentIDs[4]},
		{"结构力学", parentIDs[4]},
		// 外国语学院
		{"英美文学", parentIDs[5]},
		{"翻译", parentIDs[5]},
		{"语言学", parentIDs[5]},
		// 理学院
		{"数学分析", parentIDs[6]},
		{"高等代数", parentIDs[6]},
		{"大学物理", parentIDs[6]},
		// 材料与化工学院
		{"材料科学基础", parentIDs[7]},
		{"物理化学", parentIDs[7]},
		// 艺术与传媒学院
		{"设计基础", parentIDs[8]},
		{"数字媒体", parentIDs[8]},
		// 马克思主义学院
		{"思政类公共课", parentIDs[9]},
		// 公共课（跨专业）
		{"高等数学", parentIDs[10]},
		{"大学英语", parentIDs[10]},
		{"线性代数", parentIDs[10]},
		// 考研专区 ✦
		{"考研-数学", parentIDs[11]},
		{"考研-英语", parentIDs[11]},
		{"考研-政治", parentIDs[11]},
		{"考研-专业课", parentIDs[11]},
		// 其他
		{"考证教材", parentIDs[12]},
		{"课外书", parentIDs[12]},
	}

	for _, s := range subCats {
		DB.Create(&model.Category{Name: s.name, ParentID: s.parentID})
	}

	log.Printf("✅ 初始分类已创建（%d 个一级 + %d 个子分类）", len(cats), len(subCats))
}

func seedBanners() {
	var count int64
	DB.Model(&model.Banner{}).Count(&count)
	if count > 0 {
		return
	}

	banners := []model.Banner{
		{Title: "让知识循环，让校园更环保", ImageURL: "", SortOrder: 1, IsActive: 1},
		{Title: "低价淘到学长学姐的宝藏教材", ImageURL: "", SortOrder: 2, IsActive: 1},
		{Title: "考研资料、真题、笔记一站式淘", ImageURL: "", SortOrder: 3, IsActive: 1},
	}

	for _, b := range banners {
		DB.Create(&b)
	}
	log.Println("✅ 初始轮播图已创建（3 张）")
}
