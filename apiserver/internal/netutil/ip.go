package netutil

import (
	"os/exec"
	"strings"
)

func GetHostIP() (string, error) {
	cmd := exec.Command("hostname", "-i")
	ipAddr, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ipAddr)), nil
}
