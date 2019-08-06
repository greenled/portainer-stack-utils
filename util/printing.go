package util

import (
	"log"

	"github.com/spf13/viper"
)

func PrintVerbose(a ...interface{}) {
	if viper.GetBool("verbose") {
		log.Println(a)
	}
}

func PrintDebug(a ...interface{}) {
	if viper.GetBool("debug") {
		log.Println(a)
	}
}
