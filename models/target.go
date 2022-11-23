package models

import (
	"errors"
	"net"

	"github.com/sjlleo/traceSysBackend/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Target struct {
	gorm.Model
	TargetIP      string         `gorm:"type:varchar(60); comment:'需监测的 IP';" json:"ip"`
	TargetPort    int            `gorm:"type:int; comment:'被监测的端口'" json:"port"`
	Method        int            `gorm:"type:int; comment:'监测方法'" json:"method"`
	Interval      int            `gorm:"type:int; comment: '监测时间间隔'" json:"interval"`
	CreatedUserID int            `gorm:"type:int; comment: '创建监测规则的用户 ID'"`
	NodesID       datatypes.JSON `gorm:"type:string; comment: '监测所对应的节点 ID'" json:"nodeid"`
}

func (t *Target) TableName() string {
	return "target"
}

type TargetRes struct {
	TargetIP      string `json:"ip"`
	TargetPort    int    `json:"port"`
	Method        int    `json:"method"`
	Interval      int    `json:"interval"`
	CreatedUserID int    `json:"userid"`
	NodesID       string `json:"nodeid"`
}

func ListTargets(p *PaginationQ) error {
	var t []Target
	db := database.GetDB()
	tx := db.Model(&Target{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("TargetIP like ?", "%"+p.Parm+"%")
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

func DelTarget(id int) error {
	db := database.GetDB()
	res := db.Delete(&Target{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("监控 ID 未找到")
	}
	return nil
}

func AddTarget(t Target) error {
	db := database.GetDB()
	if addr := net.ParseIP(t.TargetIP); addr == nil {
		return errors.New("IP 格式错误")
	}

	if err := db.Create(&t).Error; err != nil {
		return err
	}

	return nil
}

func ModifyTarget(t *Target) error {
	db := database.GetDB()
	res := db.Model(&t).Where("id = ?", t.ID).Updates(&t)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}
