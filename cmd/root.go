package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "psu",
	Short:   "A CLI client for Portainer",
	Version: "is set on common/version.CurrentVersion",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.SetVersionTemplate("{{ version }}\n")
	cobra.AddTemplateFunc("version", common.BuildVersionString)

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file. (default \"$HOME/.psu.yaml)\"")
	rootCmd.PersistentFlags().StringP("log-level", "v", "info", "Log level. One of trace, debug, info, warning, error, fatal or panic.")
	rootCmd.PersistentFlags().BoolP("insecure", "i", false, "Skip Portainer SSL certificate verification.")
	rootCmd.PersistentFlags().StringP("url", "l", "", "Portainer url.")
	rootCmd.PersistentFlags().StringP("user", "u", "", "Portainer user.")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Portainer password.")
	rootCmd.PersistentFlags().StringP("auth-token", "A", "", "Portainer auth token.")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 0, "Waiting time before aborting (like 100ms, 30s, 1h20m).")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("insecure", rootCmd.PersistentFlags().Lookup("insecure"))
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("auth-token", rootCmd.PersistentFlags().Lookup("auth-token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".psu" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".psu")
	}

	viper.SetEnvPrefix("PSU")
	viper.AutomaticEnv() // read in environment variables that match

	replacer := strings.NewReplacer("-", "_", ".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

func initLogger() {
	logLevelName := viper.GetString("log-level")
	logLevel, err := logrus.ParseLevel(logLevelName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"level": logLevelName,
		}).Error("Unknown log level. Default \"info\" level will be used instead.")
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetLevel(logLevel)
}
