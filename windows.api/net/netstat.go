package net

import (
	"os/exec"
)

func NetWorkStatus() bool {
	cmd := exec.Command("ping", "baidu.com",  "-w", "3")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}