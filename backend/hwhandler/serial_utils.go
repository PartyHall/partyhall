package hwhandler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"go.bug.st/serial"
)

func readLine(port serial.Port) (string, error) {
	str := ""

	// We read char by char
	buff := make([]byte, 1)
	for {
		n, err := port.Read(buff)
		if err != nil {
			return str, err
		}

		if n == 0 {
			fmt.Println("No clue what to do here")
			return "", errors.New("byte = 0")
		}

		str += string(buff[:n])
		if strings.HasSuffix(str, "\n") || strings.HasSuffix(str, "\r") {
			break
		}
	}

	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")

	return str, nil
}

func debugMsg(msg string, args ...string) {
	hx := hex.EncodeToString([]byte(msg))

	fmt.Printf("Unknown message: %v\n\t=> 0x%v", msg, hx)
	if len(args) > 0 {
		fmt.Print("\n\tArgs: ")
		fmt.Println(strings.Join(args, " "))
	} else {
		fmt.Print(" (No args)\n")
	}
}
