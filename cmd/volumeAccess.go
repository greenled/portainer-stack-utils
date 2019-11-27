package cmd

import (
	"github.com/greenled/portainer-stack-utils/client"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/viper"
)

func init() {
	volumeAccessCmd := common.NewAccessCmd(client.ResourceVolume, "volumeName")

	volumeCmd.AddCommand(volumeAccessCmd)

	volumeAccessCmd.Flags().String("endpoint", "", "Endpoint name.")
	volumeAccessCmd.Flags().Bool("admins", false, "Permit access to this volume to administrators only.")
	volumeAccessCmd.Flags().Bool("private", false, "Permit access to this volume to the current user only.")
	volumeAccessCmd.Flags().Bool("public", false, "Permit access to this volume to any user.")
	viper.BindPFlag("volume.access.endpoint", volumeAccessCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("volume.access.admins", volumeAccessCmd.Flags().Lookup("admins"))
	viper.BindPFlag("volume.access.private", volumeAccessCmd.Flags().Lookup("private"))
	viper.BindPFlag("volume.access.public", volumeAccessCmd.Flags().Lookup("public"))
}
