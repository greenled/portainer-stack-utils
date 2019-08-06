package util

import (
	"github.com/sirupsen/logrus"
)

// CheckError checks if an error occurred (it's not nil)
func CheckError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
