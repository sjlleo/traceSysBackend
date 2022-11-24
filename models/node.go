package models

import (
	"database/sql"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/sjlleo/traceSysBackend/database"
)

type Nodes struct {
	ID            uint         `gorm:"primarykey"`
	CreatedAt     time.Time    `json:"-"`
	UpdatedAt     time.Time    `json:"-"`
	DeletedAt     sql.NullTime `gorm:"index" json:"-"`
	IP            string       `gorm:"type:varchar(50);uniqueIndex" json:"IP"`
	Role          int          `gorm:"type:int"`
	CreatedUserID uint         `gorm:"type:int"`
	Secret        string       `gorm:"type:varchar(50)" json:"secret"`
	Lastseen      time.Time
}

func (t *Nodes) TableName() string { return "nodes" }

// 分页
func (a *Admin) ListNodes(p *PaginationQ) error {
	// log.Println("管理员请求的")
	var nodes []Nodes
	db := database.GetDB()
	tx := db.Model(&Nodes{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("ip like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &nodes)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = nodes
		return nil
	}
}

func (n *Normal) ListNodes(p *PaginationQ) error {
	// log.Println("普通用户请求的")
	var nodes []Nodes
	db := database.GetDB()
	tx := db.Model(&Nodes{}).Where("created_user_id = ?", n.UserID)
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("ip like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &nodes)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = nodes
		return nil
	}
}

type NodeUser struct {
	ID uint   `json:"value"`
	IP string `json:"label"`
}

func ListNodesUser(role uint) (*[]NodeUser, error) {
	var nodes []NodeUser
	db := database.GetDB()

	tx := db.Model(&Nodes{})
	// 查找条件
	if role != 1 {
		tx = tx.Where("role=?", role)
	}
	tx.Where("deleted_at is null").Find(&nodes)
	return &nodes, nil
}

func (a *Admin) DelNode(id int) error {
	db := database.GetDB()
	res := db.Delete(&Nodes{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("节点 ID 未找到")
	}
	return nil
}

func (n *Normal) DelNode(id int) error {
	db := database.GetDB()
	res := db.Where("created_user_id = ?", n.UserID).Delete(&Nodes{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("节点 ID 未找到")
	}
	return nil
}

func (a *Admin) AddNode(ip string, role_str string, secret string) error {
	db := database.GetDB()

	if addr := net.ParseIP(ip); addr == nil {
		return errors.New("IP 格式错误")
	}

	role, err := strconv.Atoi(role_str)
	if err != nil {
		return errors.New("权限格式错误")
	}

	node := Nodes{
		IP:            ip,
		Role:          role,
		Secret:        secret,
		CreatedUserID: a.UserID,
	}

	if err := db.Create(&node).Error; err != nil {
		return err
	}

	return nil
}

func (n *Normal) AddNode(ip string, role_str string, secret string) error {
	db := database.GetDB()

	if addr := net.ParseIP(ip); addr == nil {
		return errors.New("IP 格式错误")
	}

	role, err := strconv.Atoi(role_str)
	if err != nil {
		return errors.New("权限格式错误")
	}

	if role == 1 {
		return errors.New("无权将节点设置为管理员专用")
	}

	node := Nodes{
		IP:            ip,
		Role:          role,
		Secret:        secret,
		CreatedUserID: n.UserID,
	}

	if err := db.Create(&node).Error; err != nil {
		return err
	}

	return nil
}

func (a *Admin) ModifyNode(n *Nodes) error {
	db := database.GetDB()
	res := db.Model(&Nodes{ID: n.ID}).Updates(n)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}

func (m *Normal) ModifyNode(n *Nodes) error {
	if n.Role == 1 {
		return errors.New("无权设置为管理员权限")
	}
	db := database.GetDB()
	res := db.Model(&Nodes{}).Where("created_user_id = ?", m.UserID).Where("id = ?", n.ID).Updates(n)
	if res.Error != nil {
		return res.Error
	} else {
		if res.RowsAffected == 0 {
			return errors.New("无权限修改此节点")
		}
		return nil
	}
}
