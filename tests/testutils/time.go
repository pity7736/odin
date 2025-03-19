package testutils

import "time"

func IsTimeClose(first time.Time, second time.Time) bool {
	return first.Second() == second.Second() &&
		first.Minute() == second.Minute() &&
		first.Hour() == second.Hour() &&
		first.Day() == second.Day() &&
		first.Month() == second.Month() &&
		first.Year() == second.Year()
}
