package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// List 获取分类列表（树形结构）
func (h *CategoryHandler) List(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	// 构建树形结构
	tree := buildCategoryTree(categories, 0)
	common.Success(c, tree)
}

type CategoryNode struct {
	model.Category
	Children []CategoryNode `json:"children"`
}

func buildCategoryTree(categories []model.Category, parentID uint) []CategoryNode {
	return doBuildCategoryTree(categories, parentID, 0)
}

func doBuildCategoryTree(categories []model.Category, parentID uint, depth int) []CategoryNode {
	if depth > 5 { // 限制最大深度，防循环
		return nil
	}
	var nodes []CategoryNode
	for _, cat := range categories {
		if cat.ParentID == parentID {
			node := CategoryNode{
				Category: cat,
				Children: doBuildCategoryTree(categories, cat.ID, depth+1),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}
