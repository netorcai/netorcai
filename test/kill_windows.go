// +build windows

package test

import (
	"os/exec"
)

func killallNetorcai() error {
	cmd := exec.Command("taskkill")
	cmd.Args = []string{"taskkill", "/IM", "netorcai.exe"}
	return cmd.Run()
}

func killallNetorcaiSIGKILL() error {
	cmd := exec.Command("taskkill")
	cmd.Args = []string{"taskkill", "/F", "/IM", "netorcai.exe"}
	return cmd.Run()
}
