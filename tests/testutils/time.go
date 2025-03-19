package testutils

import "time"

func IsTimeClose(first time.Time, second time.Time) bool {
	return first.Year() == second.Year() &&
		first.Month() == second.Month() &&
		first.Day() == second.Day() &&
		first.Hour() == second.Hour() &&
		first.Minute() == second.Minute() &&
		first.Second() == second.Second()
}
