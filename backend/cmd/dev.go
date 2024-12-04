package cmd

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/log"
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
