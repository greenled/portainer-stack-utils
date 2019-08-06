package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

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

func NewTabWriter(headers []string) (*tabwriter.Writer, error) {
	writer := tabwriter.NewWriter(os.Stdout, 20, 2, 3, ' ', 0)
	_, err := fmt.Fprintln(writer, strings.Join(headers, "\t"))
	if err != nil {
		return &tabwriter.Writer{}, err
	}
	return writer, nil
}
