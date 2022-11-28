package models

import (
	"log"
	"testing"

	"github.com/sjlleo/traceSysBackend/database"
)

func TestGetNodeId(t *testing.T) {
	// GetTraceList("NpBfZecKWyiv81XtM-F9R")
	// ListNodesUser(2)
}

func initDB() {
	database.InitDBConn()
	DBAutoMigration()
}

func TestExceed(t *testing.T) {
	initDB()
	task := Tasks{
		Type:      1,
		TraceType: 1,
		TTL:       12,
		NodeID:    24,
		TargetID:  16,
		ExceedRTT: 240,
	}
	res := CheckExceed(&task)
	log.Println(res)
}
