package main

import (
	"time"
)

func isFirstDayOfYear(d time.Time) bool {
	return d.Month() == 1 && d.Day() == 1
}

func isFirstDayOfMonth(d time.Time) bool {
	return d.Day() == 1
}

func isLastDayOfYear(d time.Time) bool {
	return d.Month() == 12 && d.Day() == 31
}

func isLastDayOfMonth(d time.Time) bool {
	firstOfMonth := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return d.Day() == lastOfMonth.Day()
}

func isLastDayOfWeek(d time.Time) bool {
	return d.Weekday() == 6
}

func fileAge(d time.Time) int {
	return int(time.Since(d).Hours() / 24)
}
