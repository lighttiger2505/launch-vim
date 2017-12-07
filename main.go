package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	exitStatus := launchVim()
	os.Exit(exitStatus)
}

func launchVim() int {
	message := `# |<----  Opened the file with your favorite editor. The first block of text is the title.  ---->|


# |<----  The following blocks are explanations  ---->|

`

	// Make temp editing file
	fPath := getFilePath("ISSUE")
	err := makeFile(fPath, message)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed make edit file. %s\n", err.Error()))
		return 1
	}
	defer deleteFile(fPath)

	// Open text editor
	err = openEditor("vim", "--cmd", "set ft=gitcommit tw=0 wrap lbr", fPath)
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

	// Parce read content
	reader := bytes.NewReader(content)
	title, body, err := perseTitleAndBody(reader, "#")
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed parce content. %s\n", err.Error()))
		return 1
	}

	fmt.Fprint(os.Stdout, fmt.Sprintf("title=%s, body=%s\n", title, body))

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

func makeFile(fPath, message string) (err error) {
	// only write message if file doesn't exist
	if !isFileExist(fPath) && message != "" {
		err = ioutil.WriteFile(fPath, []byte(message), 0644)
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

func perseTitleAndBody(reader io.Reader, cs string) (title, body string, err error) {
	var titleParts, bodyParts []string

	r := regexp.MustCompile("\\S")
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, cs) {
			continue
		}

		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return
}
