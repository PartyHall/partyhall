package cmd

import (
	"os"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/spf13/cobra"
)

var convertVideoCmd = &cobra.Command{
	Use:   "convert-video",
	Short: "Convert a video to the best format for the pi",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}
		// services.GET.VideoConverter.ConvertForRaspberryPi(args[0])
	},
}
