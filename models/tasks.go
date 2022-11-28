package models

import (
	"errors"

	"github.com/sjlleo/traceSysBackend/database"
	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	Name             string  `gorm:"type:varchar(255); comment: '任务名称'" json:"name"`
	Type             uint    `gorm:"type:int; comment: '任务类型'" json:"type"`
	TraceType        uint    `gorm:"type:int; comment: '路由跟踪类型'" json:"traceType"`
	CallMethod       uint    `gorm:"type:int; comment: '送信方法'" json:"callMethod"`
	TTL              uint    `gorm:"type:int; comment: '参考 TTL'" json:"ttl"`
	NodeID           uint    `gorm:"type:int; comment: '节点 ID'" json:"node_id"`
	TargetID         uint    `gorm:"type:int; comment: '监测目标 ID'" json:"targetID"`
	ExceedRTT        float64 `gorm:"type:float; comment: 'RTT 报警阈值'" json:"exceedRTT"`
	ExceedPacketLoss float64 `gorm:"type:float; comment: '丢包报警阈值'" json:"exceedPacketLoss"`
	CreatedUserID    uint    `gorm:"type:int; comment: '创建用户 ID'" json:"userid"`
}

func (m *Tasks) TableName() string {
	return "tasks"
}

func CreateTask(t *Tasks) error {
	db := database.GetDB()
	tx := db.Begin()
	err := tx.Create(t).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (A *Admin) UpdateTask(t *Tasks) error {
	db := database.GetDB()
	tx := db.Begin()
	err := tx.Model(t).Where("id =?", t.ID).Updates(&t).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (n *Normal) UpdateTask(t *Tasks) error {
	db := database.GetDB()
	tx := db.Begin()
	err := tx.Model(t).Where("id =?", t.ID).Updates(&t).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (A *Admin) DeleteTask(id uint) error {
	db := database.GetDB()
	res := db.Delete(&Target{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("任务 ID 未找到")
	}
	return nil
}

func (n *Normal) DeleteTask(id uint) error {
	db := database.GetDB()
	res := db.Where("created_user_id = ?", n.UserID).Delete(&Target{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("任务 ID 未找到")
	}
	return nil
}

func (A *Admin) GetTask(p *PaginationQ) error {
	var t []Tasks
	db := database.GetDB()
	tx := db.Model(&Target{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("name like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &t)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = t
		return nil
	}
}

func (n *Normal) GetTask(p *PaginationQ) error {
	var t []Tasks
	db := database.GetDB()
	tx := db.Model(&Target{}).Where("created_user_id = ?", n.UserID)
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("name like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &t)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = t
		return nil
	}
}

func GetAllTasks() ([]Tasks, error) {
	var t []Tasks
	db := database.GetDB()
    err := db.Model(&Tasks{}).Find(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}