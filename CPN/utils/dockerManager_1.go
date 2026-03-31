//go:build ignore

package utils

import (
	"os/exec"
)

func DockerManager(cmdString string) {
	cmd := exec.Command(cmdString)
	if err := cmd.Start(); err != nil { // 运行命令
		Logger.Errorf("Cmd: %s\nError: %s", cmdString, err)
	}
}
