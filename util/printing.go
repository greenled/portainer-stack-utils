package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func PrintDebugRequest(title string, req *http.Request) error {
	if viper.GetBool("debug") {
		var bodyString string
		if req.Body != nil {
			bodyBytes, err := ioutil.ReadAll(req.Body)
			defer req.Body.Close()
			if err != nil {
				return err
			}
			bodyString = string(bodyBytes)
			req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		}

		PrintDebug(fmt.Sprintf(`%s
---
Method: %s
URL: %s
Body:
%s
---`, title, req.Method, req.URL.String(), string(bodyString)))
	}

	return nil
}

func PrintDebugResponse(title string, resp *http.Response) error {
	if viper.GetBool("debug") {
		var bodyString string
		if resp.Body != nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return err
			}
			bodyString = string(bodyBytes)
			resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		}

		PrintDebug(fmt.Sprintf(`%s
---
Status: %s
Body:
%s
---`, title, resp.Status, bodyString))
	}

	return nil
}
