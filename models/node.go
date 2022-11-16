package models

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sjlleo/traceSysBackend/database"
)

type Nodes struct {
	ID        uint         `gorm:"primarykey"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	IP        string       `gorm:"type:varchar(50);uniqueIndex" json:"IP"`
	Role      int          `gorm:"type:int"`
	Secret    string       `gorm:"type:varchar(50)" json:"secret"`
	Lastseen  time.Time
}

func (t *Nodes) TableName() string { return "nodes" }

// 分页
func ListNodes(p *PaginationQ) error {
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

func ListNodesUser(role int) (*[]Nodes, error) {
	var nodes []Nodes
	db := database.GetDB()

	tx := db.Model(&Nodes{})
	// 查找条件
	if role != 2 {
		tx = tx.Where("role=?", role)
	}
	tx.Where("deleted_at is null").Find(&nodes)
	return &nodes, nil
}

func DelNode(id int) error {
	db := database.GetDB()
	res := db.Delete(&Nodes{}, id)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("节点 ID 未找到")
	}
	return nil
}

func AddNode(ip string, role_str string, secret string) error {
	db := database.GetDB()

	if addr := net.ParseIP(ip); addr == nil {
		return errors.New("IP 格式错误")
	}

	role, err := strconv.Atoi(role_str)
	if err != nil {
		return errors.New("权限格式错误")
	}

	node := Nodes{
		IP:     ip,
		Role:   role,
		Secret: secret,
	}

	if err := db.Create(&node).Error; err != nil {
		return err
	}

	return nil
}

func ModifyNode(n *Nodes) error {
	db := database.GetDB()
	log.Println(n)
	res := db.Model(&Nodes{ID: n.ID}).Updates(n)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}
