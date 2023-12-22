package timehelper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetZeroTime(timestamp int64) time.Time {
	t := time.Unix(timestamp, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ToTodayTime(sTime string) (time.Time, error) {
	now := time.Now()
	if len(sTime) == 0 {
		return now, errors.New("empty string")
	}
	var vec []int
	for _, fs := range strings.Split(sTime, ":") {
		if f, err := strconv.Atoi(fs); err != nil {
			return now, err
		} else {
			vec = append(vec, f)
		}
	}

	var h, m, s int
	switch len(vec) {
	case 1:
		h = vec[0]
	case 2:
		h, m = vec[0], vec[1]
	case 3:
		h, m, s = vec[0], vec[1], vec[2]
	}

	return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location()), nil
}

func StrToTime(sTime string) (time.Time, error) {
	var year, month, day, hour, minute, second int
	_, err := fmt.Sscanf(sTime, "%d-%d-%d %d:%d:%d", &year, &month, &day, &hour, &minute, &second)
	if nil != err {
		return time.Now(), err
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local), nil
}

func TimeToStr(nTime uint32) string {
	t := int64(nTime)
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// GetDaysZeroTime 获取几天后的零点时间戳
func GetDaysZeroTime(day int) time.Time {
	t := time.Now().AddDate(0, 0, day)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
