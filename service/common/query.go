package common

import (
	"fast-gin/global"
	"fast-gin/models"
	"fmt"
	"gorm.io/gorm"
)

type QueryOption struct {
	models.PageInfo
	Likes    []string
	Where    *gorm.DB
	Preloads []string
	Debug    bool
}

func QueryList[T any](model T, option QueryOption) (list []T, count int64, err error) {
	list = make([]T, 0)

	query := global.DB.Where(model)

	// Fuzzy
	if option.Key != "" {
		if len(option.Likes) != 0 {
			likeQuery := global.DB.Where("")
			for _, col := range option.Likes {
				likeQuery.Or(fmt.Sprintf("%s like ?", col), fmt.Sprintf("%%%s%%", option.Key))
			}
			query.Where(likeQuery)
		}
	}

	// Preload
	for _, preload := range option.Preloads {
		query = query.Preload(preload)
	}

	// Pagination
	if option.Page <= 0 {
		option.Page = 1
	}
	if option.Limit <= 0 {
		option.Limit = -1 // No pagination
	}
	// Pagination, offset = (page-1) * limit
	offset := (option.Page - 1) * option.Limit
	if option.Order == "" {
		option.Order = "created_at desc"
	}

	db := global.DB.Where("")
	if option.Debug {
		db = db.Debug()
	}

	db.Where(query).Limit(option.Limit).Offset(offset).Order(option.Order).Find(&list)
	db.Model(model).Where(query).Count(&count)
	return
}
