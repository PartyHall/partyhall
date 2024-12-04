package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/partyhall/partyhall/hwhandler"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var hwHandlerCmd = &cobra.Command{
	Use:   "hwhandler",
	Short: "Bridges the appliance software to the Arduino",
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

			hh := hwhandler.HardwareHandler{
				PortName:            p,
				LastButtonPress:     "",
				LastPing:            time.Now(),
				LastButtonPressTime: time.Now(),
			}
			go hh.ProcessSerialConn()
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
