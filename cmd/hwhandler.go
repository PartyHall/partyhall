package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/partyhall/partyhall/hwhandler"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var hwhandlerCmd = &cobra.Command{
	Use:   "hwhandler",
	Short: "Hardware Handler",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := hwhandler.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ports, err := serial.GetPortsList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(ports) == 0 {
			fmt.Println("No arduino found")
			os.Exit(1)
		}

		for _, p := range ports {
			if !strings.HasPrefix(p, "/dev/ttyUSB") {
				continue
			}

			hh := hwhandler.HardwareHandler{PortName: p}
			go hh.ProcessSerialConn()
		}

		c := make(chan os.Signal)
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
