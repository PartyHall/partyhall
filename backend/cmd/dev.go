package cmd

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/pipewire"
	"github.com/partyhall/partyhall/utils"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development related commands",
}

var generateJwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Generate a JWT for the appliance",
	Run: func(cmd *cobra.Command, args []string) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.GetClaimsFromAppliance())
		tokenString, err := token.SignedString(config.GET.Mercure.SubscriberKey)
		if err != nil {
			log.LOG.Errorw("Failed to generate JWT", "err", err)

			return
		}

		fmt.Println(tokenString)
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Some test in dev, should not be used in prod",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := pipewire.GetDevices()
		if err != nil {
			fmt.Println(err)
			return
		}

		var src *models.PwDevice
		fmt.Println("Sources:")
		for _, s := range devices.Sources {
			fmt.Println(s)
			if s.ID == 67 {
				src = &s
			}
		}

		fmt.Println()

		var dst *models.PwDevice
		fmt.Println("Sinks:")
		for _, s := range devices.Sinks {
			fmt.Println(s)
			if s.ID == 53 {
				dst = &s
			}
		}

		fmt.Println()

		fmt.Println("Links:")
		for _, s := range devices.Links {
			fmt.Println(s)
		}

		err = pipewire.LinkDevice(src, dst)
		if err != nil {
			fmt.Println(err)
		}

		/*
			fmt.Println(pipewire.SetVolume(dev, 0.8), *dev)
			time.Sleep(5 * time.Second)
			fmt.Println(pipewire.SetVolume(dev, 1.2), *dev)
			time.Sleep(5 * time.Second)
			fmt.Println(pipewire.SetVolume(dev, 1.0), *dev)
		*/
	},
}
