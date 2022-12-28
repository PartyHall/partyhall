package main

import (
	"fmt"
	"os"

	"github.com/partyhall/partyhall/cmd"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/logs"
)

func main() {
	if err := config.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logs.Init()

	cmd.Execute()
}
