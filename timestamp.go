package timehelper

import "time"

func GetZeroTime(timestamp int64) int64 {
	t := time.Unix(timestamp, 0)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Unix()
}
