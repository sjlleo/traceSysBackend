package models

import (
	"encoding/json"

	"github.com/sjlleo/traceSysBackend/database"
	"gorm.io/gorm"
)

type Result struct {
	gorm.Model
	NodeID   int     `gorm:"type:bigint(20)" json:"node_id"`
	TargetID int     `gorm:"type:bigint(20)" json:"target_id"`
	TTL      int     `gorm:"type:int" json:"ttl"`
	AvgRTT   float64 `gorm:"type:float" json:"avgRTT"`
	MaxRTT   float64 `gorm:"type:float" json:"maxRTT"`
	MinRTT   float64 `gorm:"type:float" json:"minRTT"`
	IPList   string  `gorm:"type:string" json:"ip_list"`
	Interval int     `gorm:"type:int" json:"interval"`
}

type HopReport struct {
	IPList     []string `json:"ip_list"`
	MinLatency float64  `json:"min_latency"`
	MaxLatency float64  `json:"max_latency"`
	AvgLatency float64  `json:"avg_latency"`
}

type ClientData struct {
	Data     map[int]*HopReport
	Interval int  `json:"interval"`
	NodeID   uint `json:"nodeId"`
	TaskID   uint `json:"taskId"`
}

func AddTraceData(c *ClientData) error {
	var err error
	db := database.GetDB()
	for ttl, v := range c.Data {
		if v.IPList == nil {
			continue
		}
		ipStr, _ := json.Marshal(v.IPList)
		r := Result{
			IPList:   string(ipStr),
			TTL:      ttl + 1,
			AvgRTT:   v.AvgLatency,
			MaxRTT:   v.MaxLatency,
			MinRTT:   v.MinLatency,
			Interval: c.Interval,
			NodeID:   int(c.NodeID),
			TargetID: int(c.TaskID),
		}
		err = db.Model(&Result{}).Create(&r).Error
		if err != nil {
			return err
		}
	}
	return nil
}
