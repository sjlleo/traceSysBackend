package models

import (
	"log"
	"strconv"
	"strings"

	"github.com/sjlleo/traceSysBackend/database"
)

type TraceList struct {
	Task []TraceTask `json:"task"`
}

type TraceTask struct {
	TaskID   uint   `json:"id"`
	NodeID   uint   `json:"nodeId"`
	Interval int    `json:"interval"`
	Method   int    `json:"method"`
	IP       string `json:"ip"`
}

func GetTraceList(token string) TraceList {
	var nodes Nodes
	var target []Target
	var list TraceList
	db := database.GetDB()
	// 获取 Token 对应的 node id
	db.Model(&nodes).Where("secret = ?", token).Take(&nodes)
	nodeIdStr := strconv.Itoa(int(nodes.ID))
	db.Model(&Target{}).Where("nodes_id LIKE ?", "%"+nodeIdStr+"%").Find(&target)
	log.Println(target)
	for _, v := range target {
		split_arr := strings.Split(v.NodesID, ",")
		// Python: if v.nodes_id in split_arr
		for _, k := range split_arr {
			if k == nodeIdStr {
				// 将数据组成 task
				task := TraceTask{
					TaskID:   v.ID,
					NodeID:   nodes.ID,
					Interval: v.Interval,
					Method:   v.Method,
					IP:       v.TargetIP,
				}
				// 将 task 放入 list
				list.Task = append(list.Task, task)
			}
		}
	}
	return list
}
