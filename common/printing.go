package common

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func NewTabWriter(headers []string) (*tabwriter.Writer, error) {
	writer := tabwriter.NewWriter(os.Stdout, 20, 2, 3, ' ', 0)
	_, err := fmt.Fprintln(writer, strings.Join(headers, "\t"))
	if err != nil {
		return &tabwriter.Writer{}, err
	}
	return writer, nil
}
