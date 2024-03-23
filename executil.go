package executil

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/anmitsu/go-shlex"
)

func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func Command(cmd string) (*exec.Cmd, error) {
	// When GOOS == "windows", we need to use non-POSIX rules for splitting the command.
	// This ensures compatibility with Windows shell syntax.
	args, err := shlex.Split(cmd, runtime.GOOS != "windows")
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
