package utils

import (
	"time"
)

// GetStartAndEndOfDay returns the Unix timestamps for the start and end of the current day
func GetStartAndEndOfDay() (int64, int64) {
	location := time.Local

	now := time.Now().In(location)

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	startTime := startOfDay.Unix()
	endTime := now.Unix()

	return startTime, endTime
}
