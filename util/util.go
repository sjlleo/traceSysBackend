package util

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"time"
)

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func IN_10_Minutes(t1, t2 time.Time) bool {
	// log.Println(math.Floor(t1.Sub(t2).Minutes()))
	return int(math.Floor(t1.Sub(t2).Minutes()+10))%1440 < 10
}
