package utils

import (
	"fmt"
	"time"
)

func GetTimeFromInt(intTime int64) time.Time {
	return time.Unix(intTime, 0)
}

func GetHourAsString(time time.Time) string {
	return fmt.Sprintf("%02d", time.Hour())
}

func GetWeekdayAsString(time time.Time) string {
	return time.Weekday().String()
}
