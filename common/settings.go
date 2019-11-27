package common

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// LoadSettings loads the settings file currently used by viper into a new viper instance
func LoadSettings() (v *viper.Viper, err error) {
	// Set settings file name
	var settingsFile string
	if viper.ConfigFileUsed() != "" {
		// Use settings file from viper
		settingsFile = viper.ConfigFileUsed()
	} else {
		// Find home directory
		var home string
		home, err = homedir.Dir()
		if err != nil {
			return
		}

		// Use $HOME/.psu.yaml
		settingsFile = fmt.Sprintf("%s%s.psu.yaml", home, string(os.PathSeparator))
	}
	v = viper.New()
	v.SetConfigFile(settingsFile)

	// Read settings from file
	err = v.ReadInConfig()

	return
}

// CheckSettingKeyExists checks a given setting key exists in the default viper
func CheckSettingKeyExists(key string) (keyExists bool) {
	for _, k := range viper.AllKeys() {
		if k == key {
			keyExists = true
			break
		}
	}
	return
}
