package util

import (
	"fmt"
	"log"
)

// CheckError checks if an error occurred (it's not nil)
func CheckError(err error) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error: %s", err.Error()))
	}
}
