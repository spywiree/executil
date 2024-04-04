package executil

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"mvdan.cc/sh/v3/shell"
)

func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func Command(cmd string) (*exec.Cmd, error) {
	args, err := shell.Fields(cmd, nil)
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		return nil, err
	}

	if !IsFileExists(path) {
		path = args[0]
	}

	return exec.Command(path, args[1:]...), nil
}

func RedirectIO(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func SetWdToParent(cmd *exec.Cmd) {
	cmd.Dir = filepath.Dir(cmd.Args[0])
}
