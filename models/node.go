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
	ID        uint         `gorm:"primarykey"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	IP        string       `gorm:"type:varchar(50);uniqueIndex" json:"IP,omitempty"`
	Role      int          `gorm:"type:int"`
	Secret    string       `gorm:"type:varchar(50)" json:"secret,omitempty"`
	Lastseen  time.Time
}

func (t *Nodes) TableName() string { return "nodes" }

func ListNodes() (*[]Nodes, error) {
	var nodes []Nodes
	db := database.GetDB()
	res := db.Find(&nodes)
	if res.Error != nil {
		return nil, res.Error
	} else {
		return &nodes, nil
	}
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
