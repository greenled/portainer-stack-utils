package cmd

import (
	"io/ioutil"
	"net/http"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/greenled/portainer-stack-utils/common"

	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start an HTTP proxy to an endpoint's Docker API",
	Example: `  Expose "primary" endpoint's Docker API at 127.0.0.1:11000: 
  psu proxy --endpoint primary --address 127.0.0.1:11000

  Configure local Docker client to connect to proxy:
  export DOCKER_HOST=tcp://127.0.0.1:11000`,
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, err := common.GetClient()
		common.CheckError(err)

		var endpoint portainer.Endpoint
		if endpointName := viper.GetString("proxy.endpoint"); endpointName == "" {
			// Guess endpoint if not set
			logrus.WithFields(logrus.Fields{
				"implications": "Command will fail if there is not exactly one endpoint available",
			}).Warning("Endpoint not set")
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetDefaultEndpoint()
			common.CheckError(endpointRetrievalErr)
			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint.Name,
			}).Debug("Using the only available endpoint")
		} else {
			// Get endpoint by name
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetEndpointByName(endpointName)
			common.CheckError(endpointRetrievalErr)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			resp, err := portainerClient.Proxy(endpoint.ID, r)
			if err != nil {
				logrus.Fatal(err)
			}

			for key, value := range resp.Header {
				for i := range value {
					w.Header().Add(key, value[i])
				}
			}

			w.WriteHeader(resp.StatusCode)

			if resp.Body != nil {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					logrus.Fatal(err)
				}
				w.Write(bodyBytes)
			}
		})

		err = http.ListenAndServe(viper.GetString("proxy.address"), nil)
		if err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().String("endpoint", "", "Endpoint name.")
	proxyCmd.Flags().String("address", "127.0.0.1:2375", "Address to bind to.")
	viper.BindPFlag("proxy.endpoint", proxyCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("proxy.address", proxyCmd.Flags().Lookup("address"))
}
