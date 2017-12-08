package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	exitStatus := launchVim()
	os.Exit(exitStatus)
}

func launchVim() int {
	// Make temp editing file
	fPath := getFilePath("ISSUE")
	err := makeFile(fPath)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed make edit file. %s\n", err.Error()))
		return 1
	}
	defer deleteFile(fPath)

	// Open text editor
	err = openEditor("vim", fPath)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed open text editor. %s\n", err.Error()))
		return 1
	}

	// Read edit file
	content, err := ioutil.ReadFile(fPath)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed read content. %s\n", err.Error()))
		return 1
	}

	fmt.Fprint(os.Stdout, string(content))

	return 0
}

func getFilePath(about string) string {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")
	}
	fname := filepath.Join(home, "tmp", fmt.Sprintf("%s_EDITMSG", about))
	return fname
}

func makeFile(fPath string) (err error) {
	if !isFileExist(fPath) {
		err = ioutil.WriteFile(fPath, []byte(""), 0644)
		if err != nil {
			return
		}
	}
	return
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func deleteFile(fPath string) error {
	return os.Remove(fPath)
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
