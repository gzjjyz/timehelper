package timehelper

import "time"

func TimestampSameDay(t1, t2 int64) bool {
	if t1 == t2 {
		return true
	}
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	left := time.Unix(t1, 0)
	right := time.Unix(t2, 0)

	return left.Year() == right.Year() && left.Month() == right.Month() && left.Day() == right.Day()
}

func TimestampSameWeek(t1, t2 int64) bool {
	if t1 == t2 {
		return true
	}
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	ly, lm := time.Unix(t1, 0).ISOWeek()
	ry, rm := time.Unix(t2, 0).ISOWeek()

	return ly == ry && lm == rm
}

func TimestampSubDays(t1, t2 int64) int32 {
	if t1 > t2 {
		t2, t1 = t1, t2
	}

	return int32(GetZeroTime(t2).Unix()/86400.0 - GetZeroTime(t1).Unix()/86400.0)
}
