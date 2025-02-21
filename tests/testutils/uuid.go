package testutils

import "regexp"

func IsUUIDv7(value string) bool {
	uuidPattern, _ := regexp.Compile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidPattern.MatchString(value)

}
