package models

import (
	"encoding/json"
	"time"

	"github.com/sjlleo/traceSysBackend/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Result struct {
	NodeID     int            `gorm:"type:bigint(20)" json:"node_id"`
	TargetID   int            `gorm:"type:bigint(20)" json:"target_id"`
	TTL        int            `gorm:"type:int" json:"ttl"`
	PacketLoss float64        `gorm:"type:float" json:"packet_loss"`
	AvgRTT     float64        `gorm:"type:float" json:"avgRTT"`
	MaxRTT     float64        `gorm:"type:float" json:"maxRTT"`
	MinRTT     float64        `gorm:"type:float" json:"minRTT"`
	IPList     datatypes.JSON `gorm:"type:string" json:"ip_list"`
	Interval   int            `gorm:"type:int" json:"interval"`
	Method     int            `gorm:"type:int" json:"method"`
	gorm.Model
}

type HopReport struct {
	IPList     []string `json:"ip_list"`
	PacketLoss float64  `json:"packetLoss"`
	MinLatency float64  `json:"min_latency"`
	MaxLatency float64  `json:"max_latency"`
	AvgLatency float64  `json:"avg_latency"`
}

type ClientData struct {
	Data     map[int]*HopReport
	Interval int       `json:"interval"`
	NodeID   uint      `json:"nodeId"`
	TaskID   uint      `json:"taskId"`
	Method   int       `json:"method"`
	Time     time.Time `json:"time"`
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
			IPList:     ipStr,
			TTL:        ttl + 1,
			Method:     c.Method,
			AvgRTT:     v.AvgLatency,
			MaxRTT:     v.MaxLatency,
			MinRTT:     v.MinLatency,
			PacketLoss: v.PacketLoss,
			Interval:   c.Interval,
			NodeID:     int(c.NodeID),
			TargetID:   int(c.TaskID),
		}
		r.CreatedAt = c.Time
		err = db.Model(&Result{}).Create(&r).Error
		if err != nil {
			return err
		}
	}
	// 修改在线时间
	db.Model(&Nodes{ID: c.NodeID}).Update("lastseen", time.Now())
	return nil
}

type ShowResArgs struct {
	Method    int       `json:"method"`
	IP        string    `json:"targetIP"`
	NodeID    int       `json:"nodeID"`
	StartDate LocalTime `json:"startDate"`
	EndDate   LocalTime `json:"endDate"`
}

type FrontendResult struct {
	TTL        int            `gorm:"type:int" json:"ttl"`
	PacketLoss float64        `gorm:"type:float" json:"packet_loss"`
	AvgRTT     float64        `gorm:"type:float" json:"avgRTT"`
	MaxRTT     float64        `gorm:"type:float" json:"maxRTT"`
	MinRTT     float64        `gorm:"type:float" json:"minRTT"`
	IPList     datatypes.JSON `gorm:"type:string" json:"ip_list"`
	Interval   int            `gorm:"type:int" json:"interval"`
	CreatedAt  LocalTime      `json:"created_time"`
}

func ShowTraceData(args ShowResArgs) ([]FrontendResult, error) {
	var t Target
	var r []FrontendResult
	// 搜索监控的 IP 对应的 ID 号码
	db := database.GetDB()
	db.Model(&Target{}).Where("target_ip = ?", args.IP).Take(&t)

	tx := db.Model(&Result{})
	tx = tx.Where("method = ?", args.Method)
	tx = tx.Where("target_id = ?", t.ID).Where("node_id = ?", args.NodeID)

	startDateValid := args.StartDate.String() != "0001-01-01 00:00:00"
	endDateValid := args.EndDate.String() != "0001-01-01 00:00:00"

	if startDateValid {
		tx = tx.Where("unix_timestamp(created_at) > unix_timestamp(?)", args.StartDate)
	}
	if endDateValid {
		tx = tx.Where("unix_timestamp(created_at) < unix_timestamp(?)", args.EndDate)
	}

	if !startDateValid && !endDateValid {
		tx = tx.Limit(30).Order("created_at DESC")
	}

	err := tx.Find(&r).Error
	return r, err
}
