package internal

import (
	"os"
	"os/exec"
)

const (
	EnvFinishUpdate = "FINISH_UPDATE"
)

var _ OsOperations = (*OsOperationsImpl)(nil)

type OsOperations interface {
	Restart(path string) error
}

type OsOperationsImpl struct{}

func (OsOperationsImpl) Restart(path string) error {
	env := os.Environ()
	env = append(env, EnvFinishUpdate+"=1")
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env

	//ignore an error just start concurrently
	go cmd.Run()

	// exit the current process immediately
	os.Exit(0)

	return nil
}
