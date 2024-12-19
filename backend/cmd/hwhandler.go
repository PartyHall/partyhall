package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/partyhall/partyhall/hwhandler"
	"github.com/spf13/cobra"
)

var hwHandlerCmd = &cobra.Command{
	Use:   "hwhandler",
	Short: "Bridges the appliance software to the Arduino",
	Run: func(cmd *cobra.Command, args []string) {
		handler, err := hwhandler.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = handler.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Exit(1)
		}()

		for {
			runtime.Gosched()
		}
	},
}
