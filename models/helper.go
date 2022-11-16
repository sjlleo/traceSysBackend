package models

import (
	"gorm.io/gorm"
)

// PaginationQ gin handler query binding struct
type PaginationQ struct {
	// Ok    bool        `json:"ok"`
	Size  int         `form:"size" json:"size"`
	Page  int         `form:"page" json:"page"`
	Parm  string      `form:"parm" json:"parm"`
	Data  interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"` // save pagination list
	Total uint        `json:"total"`
}

func crudAll(p *PaginationQ, queryTx *gorm.DB, list interface{}) (int64, error) {
	//1.默认参数
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}

	//2.查询数量
	var total int64
	err := queryTx.Count(&total).Error
	if err != nil {
		return 0, err
	}
	offset := p.Size * (p.Page - 1)

	//3.偏移量的数据
	err = queryTx.Limit(p.Size).Offset(offset).Find(list).Error
	if err != nil {
		return 0, err
	}

	return total, err
}
