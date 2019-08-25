package common

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// LoadCofig loads the configuration file currently used by viper into a new viper instance
func LoadCofig() (v *viper.Viper, err error) {
	// Set config file name
	var configFile string
	if viper.ConfigFileUsed() != "" {
		// Use config file from viper
		configFile = viper.ConfigFileUsed()
	} else {
		// Find home directory
		var home string
		home, err = homedir.Dir()
		if err != nil {
			return
		}

		// Use $HOME/.psu.yaml
		configFile = fmt.Sprintf("%s%s.psu.yaml", home, string(os.PathSeparator))
	}
	v = viper.New()
	v.SetConfigFile(configFile)

	// Read config from file
	err = v.ReadInConfig()

	return
}

// CheckConfigKeyExists checks a given configuration key exists in the default viper
func CheckConfigKeyExists(key string) (keyExists bool) {
	for _, k := range viper.AllKeys() {
		if k == key {
			keyExists = true
			break
		}
	}
	return
}
