package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/partyhall/partyhall/logs"
	"github.com/partyhall/partyhall/services"
	"github.com/spf13/cobra"
)

var spotifySearchCmd = &cobra.Command{
	Use:   "spotify-search",
	Short: "Convert a video to the best format for the pi",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.Load(); err != nil {
			logs.Error(err)
			os.Exit(1)
		}
		resp, err := services.GET.Spotify.SearchSong(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, t := range resp {
			fmt.Println(t)
			for _, i := range t.Album.Images {
				if i.Width == 300 && i.Height == 300 {
					fmt.Printf("\t- %v\n", i.URL)
				}
			}
		}
	},
}
