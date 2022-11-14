package internal

import (
	"os"
	"os/exec"
)

var _ OsOperations = (*OsOperationsImpl)(nil)

type OsOperations interface {
	Restart(path string) error
}

type OsOperationsImpl struct{}

func (OsOperationsImpl) Restart(path string) error {
	env := os.Environ()
	env = append(env, "FINISH_UPDATE=1")
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env
	err := cmd.Run()
	if err == nil {
		os.Exit(0)
	}
	return err
}
