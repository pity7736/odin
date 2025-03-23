package testutils

import "os"

func GetTestPath() string {
	pwd, _ := os.Getwd()
	return pwd
}
