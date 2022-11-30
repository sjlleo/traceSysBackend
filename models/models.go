package models

import (
	"database/sql/driver"
	"time"
)

type User interface {
	ListNodes(p *PaginationQ) error
	AddNode(ip string, role_str string, secret string) error
	ModifyNode(n *Nodes) error
	DelNode(id int) error
	ListTargets(p *PaginationQ) error
	DelTarget(id int) error
	ModifyTarget(t *Target) error
	FindTargetIPNodeInfo(t *NodeInfo) error
	GetTask(p *PaginationQ) error
	UpdateTask(t *Tasks) error
	DeleteTask(id uint) error
	ListTargetUser(nodeID uint) ([]TargetUser, error)
}

type Normal struct {
	UserID uint
	RoleID uint
}

type Admin struct {
	UserID uint
	RoleID uint
}

const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
