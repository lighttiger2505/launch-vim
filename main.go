package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	exitStatus := launchVim()
	os.Exit(exitStatus)
}

func launchVim() int {
	// Open text editor
	err := openEditor("vim", "--cmd", "set ft=gitcommit tw=0 wrap lbr")
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed open text editor. %s\n", err.Error()))
		return 1
	}
	return 0
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
