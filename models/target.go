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
	TargetIP      string         `gorm:"type:varchar(60); comment:'需监测的 IP'; " json:"ip"`
	TargetPort    int            `gorm:"type:int; comment:'被监测的端口'" json:"port"`
	Method        int            `gorm:"type:int; comment:'监测方法'" json:"method"`
	Interval      int            `gorm:"type:int; comment: '监测时间间隔'" json:"interval"`
	CreatedUserID uint           `gorm:"type:int; comment: '创建监测规则的用户 ID'"`
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

func (a *Admin) ListTargets(p *PaginationQ) error {
	var t []Target
	db := database.GetDB()
	tx := db.Model(&Target{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("target_ip like ?", "%"+p.Parm+"%")
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

func (n *Normal) ListTargets(p *PaginationQ) error {
	var t []Target
	db := database.GetDB()
	tx := db.Model(&Target{}).Where("created_user_id = ?", n.UserID)
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("target_ip like ?", "%"+p.Parm+"%")
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

func (a *Admin) CountTarget() (int64, error) {
	var count int64
	db := database.GetDB()
	err := db.Model(&Target{}).Count(&count).Error
	return count, err
}

func (n *Normal) CountTarget() (int64, error) {
	var count int64
	db := database.GetDB()
	err := db.Model(&Target{}).Where("created_user_id = ?", n.UserID).Error
	return count, err
}

func (a *Admin) DelTarget(id int) error {
	db := database.GetDB()
	res := db.Delete(&Target{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("监控 ID 未找到")
	}
	return nil
}

func (n *Normal) DelTarget(id int) error {
	db := database.GetDB()
	res := db.Where("created_user_id = ?", n.UserID).Delete(&Target{}, id)
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

	var tmp Target
	db.Where("target_ip = ?", t.TargetIP).Take(&tmp)
	if tmp.TargetIP != "" {
		return errors.New("监控 IP 已经添加")
	}

	if err := db.Create(&t).Error; err != nil {
		return err
	}

	return nil
}

func (a *Admin) ModifyTarget(t *Target) error {
	db := database.GetDB()
	res := db.Model(&t).Where("id = ?", t.ID).Updates(&t)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}

func (n *Normal) ModifyTarget(t *Target) error {
	db := database.GetDB()
	res := db.Model(&Target{}).Where("created_user_id = ?", n.UserID).Where("id = ?", t.ID).Updates(&t)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return errors.New("无权限修改此节点")
		}
		return nil
	}
}

type NodeInfo struct {
	Code     uint           `gorm:"-" json:"code"`
	TargetIP string         `json:"ip"`
	NodesID  datatypes.JSON `json:"nodeid"`
}

func FindTargetIPByID(targetID uint) (Target, error) {
	db := database.GetDB()
	t := Target{}
	err := db.Model(&t).Where("id = ?", targetID).Take(&t).Error
	return t, err
}

func (a *Admin) FindTargetIPNodeInfo(t *NodeInfo) error {
	db := database.GetDB()
	err := db.Model(&Target{}).Where("target_ip = ?", t.TargetIP).Take(&t).Error
	return err
}

func (n *Normal) FindTargetIPNodeInfo(t *NodeInfo) error {
	db := database.GetDB()
	err := db.Model(&Target{}).Where("target_ip = ?", t.TargetIP).Where("created_user_id = ?", n.UserID).Take(&t).Error
	return err
}

type TargetUser struct {
	ID       uint   `json:"value"`
	TargetIP string `json:"label"`
}

func (a *Admin) ListTargetUser(nodeID uint) ([]TargetUser, error) {
	var targets []TargetUser
	db := database.GetDB()

	tx := db.Model(&Target{})

	tx.Where("deleted_at is null").Where("JSON_CONTAINS(`nodes_id` ->> '$[*]',JSON_ARRAY(?),'$')", nodeID).Find(&targets)
	return targets, nil
}

func (n *Normal) ListTargetUser(nodeID uint) ([]TargetUser, error) {
	var targets []TargetUser
	db := database.GetDB()

	tx := db.Model(&Target{})

	tx.Where("deleted_at is null").Where("JSON_CONTAINS(`nodes_id` ->> '$[*]',JSON_ARRAY(?),'$')", nodeID).Where("created_user_id = ?", n.UserID).Find(&targets)
	return targets, nil
}
