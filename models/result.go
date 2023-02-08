package models

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/sjlleo/traceSysBackend/database"
	"github.com/sjlleo/traceSysBackend/util"
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

func CheckExceed(t *Tasks) bool {
	var r []Result
	db := database.GetDB()
	tx := db.Model(&Result{})
	tx = tx.Where("method = ?", t.TraceType).Where("ttl = ?", t.TTL)
	tx = tx.Where("target_id = ?", t.TargetID).Where("node_id = ?", t.NodeID)

	// 这里获取近一个月的时间作为参考数据，时间太久远的数据除了用来考古，没有太大的参考意义
	nowTime := time.Now()
	nowTimeStr := nowTime.Format("2006-01-02")
	previosTime := nowTime.AddDate(0, -1, 0)
	previosTimeStr := previosTime.Format("2006-01-02")
	// 时间相关的语句目前 Gorm 还不能很完美的处理，直接手撕 SQL 语句
	tx = tx.Where("`created_at` BETWEEN ? AND ?", previosTimeStr, nowTimeStr)

	err := tx.Scan(&r).Error

	if err != nil {
		return false
	}

	var avgMaxRTT, avgPacketLoss float64
	var count uint
	for _, v := range r {
		if len(v.IPList) == 0 {
			continue
		}
		// 遍历搜索当前时间往后的10分钟内的历史数据
		if util.IN_10_Minutes(nowTime, v.CreatedAt) {
			// 单独根据某个数据集的峰值数据判断拥塞意义不大，还容易因过度耦合产生误报
			switch t.Type {
			case 1:
				// RTT
				avgMaxRTT += v.MaxRTT
			case 2:
				// Packet Loss
				avgPacketLoss += v.PacketLoss
			}
			log.Println(v.CreatedAt)
			count++
		}
	}

	switch t.Type {
	case 1:
		log.Println(avgMaxRTT / float64(count))
		return avgMaxRTT/float64(count) > t.ExceedRTT
	case 2:
		return avgPacketLoss/float64(count) > t.ExceedPacketLoss
	default:
		return false
	}
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
	Count      int            `gorm:"-" json:"-"`
}

func ShowTraceData(args ShowResArgs) ([]FrontendResult, error) {
	var t Target
	var r []FrontendResult
	// 搜索监控的 IP 对应的 ID 号码
	db := database.GetDB()
	db.Model(&Target{}).Where("target_ip = ?", args.IP).Take(&t)

	log.Println(args)

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
		err := tx.Find(&r).Error
		log.Println(r)
		return r, err

	}

	err := tx.Find(&r).Error
	// return r, err

	if len(r) == 0 {
		return []FrontendResult{}, nil
	}

	var roundCycle int
	// 如果两个都满足，因为巨大的数据量，我们需要对数据集做一定的归并处理
	diffHour := time.Time(args.EndDate).Sub(time.Time(args.StartDate)).Hours()
	switch {
	case diffHour <= 2 || !startDateValid && !endDateValid:
		// 对于2小时以内的数据我们不做处理
		return r, err
	case diffHour <= 24:
		// 对于24小时以内的数据，我们应该将间隔调整为5分钟及以上
		if interval := r[len(r)-1].Interval; interval < 15 {
			// 设法将其调整为5分钟以上

			// 如果不能整除，则可以向上取整 5/2 = 3，变成6分钟的间隔
			roundCycle = int(math.Ceil(5 / float64(interval/3)))
			log.Println(roundCycle)
		}
	case diffHour <= 72:
		// 对于72小时以内的数据，我们应该将间隔调整为10分钟及以上
		if interval := r[len(r)-1].Interval; interval < 30 {
			// 10/2 = 5，变成10分钟的间隔
			roundCycle = int(math.Ceil(10 / float64(interval/3)))
		}
	case diffHour > 72:
		// 超过72小时的数据，间隔最好设置在20分钟及以上
		if interval := r[len(r)-1].Interval; interval < 60 {
			// 20/2 = 10，变成20分钟的间隔
			roundCycle = int(math.Ceil(20 / float64(interval/3)))
		}
	}

	// 重新定义一个新的返回结果集
	var res []FrontendResult
	var tmp [31]FrontendResult
	var count int
	// 根据计算出来的 roundCycle 开始处理数据
	for index, target := range r {
		if index != len(r)-1 && target.CreatedAt != r[index+1].CreatedAt {
			count++
			if count%roundCycle == 0 {
				for _, tmp_target := range tmp {
					if tmp_target.Interval == 0 {
						continue
					}

					tmp_target.AvgRTT /= float64(tmp_target.Count)
					tmp_target.PacketLoss /= float64(tmp_target.Count)
					tmp_target.CreatedAt = target.CreatedAt
					res = append(res, tmp_target)
				}
				tmp = [31]FrontendResult{}
			}
		} else if index == len(r)-1 && count/roundCycle > 0 {
			for _, tmp_target := range tmp {
				if tmp_target.Interval == 0 {
					continue
				}
				tmp_target.AvgRTT /= float64(count%roundCycle) + 1
				tmp_target.PacketLoss /= float64(count%roundCycle) + 1
				tmp_target.CreatedAt = target.CreatedAt
				res = append(res, tmp_target)
			}
		}
		tmp[target.TTL].TTL = target.TTL
		tmp[target.TTL].AvgRTT += target.AvgRTT
		if tmp[target.TTL].MaxRTT < target.MaxRTT {
			tmp[target.TTL].MaxRTT = target.MaxRTT
		}
		if tmp[target.TTL].MinRTT > target.MinRTT || tmp[target.TTL].MinRTT == 0 {
			tmp[target.TTL].MinRTT = target.MinRTT
		}
		tmp[target.TTL].PacketLoss += target.PacketLoss
		tmp[target.TTL].IPList = target.IPList
		tmp[target.TTL].CreatedAt = target.CreatedAt
		tmp[target.TTL].Interval = target.Interval * roundCycle
		tmp[target.TTL].Count += 1
	}
	return res, err
}
