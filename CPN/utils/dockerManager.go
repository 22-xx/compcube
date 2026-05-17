package utils

import (
	"os/exec"
)

func DockerManager(args []string) {
	go func() {
		cmd := exec.Command("docker", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			Logger.Errorf("Docker run failed: %s\nOutput: %s", err, string(output))
		}
	}()
}
