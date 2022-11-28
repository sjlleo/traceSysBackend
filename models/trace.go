package models

import (
	"github.com/sjlleo/traceSysBackend/database"
)

type TraceList struct {
	Task []TraceTask `json:"task"`
}

type TraceTask struct {
	TaskID     uint   `json:"id"`
	NodeID     uint   `json:"nodeId"`
	Interval   int    `json:"interval"`
	Method     int    `json:"method"`
	TargetPort uint   `json:"targetPort"`
	IP         string `json:"ip"`
}

func GetTraceList(token string) TraceList {
	var nodes Nodes
	var target []Target
	var list TraceList
	db := database.GetDB()
	// 获取 Token 对应的 node id
	db.Model(&nodes).Where("secret = ?", token).Take(&nodes)
	db.Raw("SELECT * FROM `target` WHERE JSON_CONTAINS(`nodes_id` ->> '$[*]',JSON_ARRAY(?),'$') AND deleted_at IS NULL", nodes.ID).Scan(&target)
	for _, v := range target {
		task := TraceTask{
			TaskID:     v.ID,
			NodeID:     nodes.ID,
			Interval:   v.Interval,
			Method:     v.Method,
			IP:         v.TargetIP,
			TargetPort: uint(v.TargetPort),
		}
		// 将 task 放入 list
		list.Task = append(list.Task, task)
	}
	return list
}
