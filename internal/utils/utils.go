package utils

import "time"

func GetTimeBoundary(timestamp int64) (start int64, end int64) {
	start = GetStartofDay(timestamp)
	end = GetEndofDay(timestamp)

	return
}

func GetStartofDay(timestamp int64) (start int64) {
	day := time.Unix(timestamp, 0)
	start = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Unix()
	return
}

func GetEndofDay(timestamp int64) (end int64) {
	day := time.Unix(timestamp, 0)
	end = time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location()).Unix()
	return
}
