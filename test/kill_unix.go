// +build !windows

package test

import (
	"os/exec"
)

func killallNetorcai() error {
	cmd := exec.Command("killall")
	cmd.Args = []string{"killall", "--quiet", "netorcai", "netorcai.cover"}
	return cmd.Run()
}

func killallNetorcaiSIGKILL() error {
	cmd := exec.Command("killall")
	cmd.Args = []string{"killall", "-KILL", "--quiet", "netorcai", "netorcai.cover"}
	return cmd.Run()
}
