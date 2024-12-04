package services

import "os/exec"

func Shutdown() error {
	err := DB.Close()
	if err != nil {
		return err
	}

	err = exec.Command("shutdown", "-h", "now").Run()

	return err
}
