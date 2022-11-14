package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	EnvFinishUpdate = "FINISH_UPDATE"
	EnvKillThisPid  = "KILL_THIS_PID"
)

var _ OsOperations = (*OsOperationsImpl)(nil)

type OsOperations interface {
	Restart(path string) error
}

type OsOperationsImpl struct{}

func (OsOperationsImpl) Restart(path string) error {
	env := os.Environ()
	env = append(env, EnvFinishUpdate+"=1")
	env = append(env, fmt.Sprintf("%v=%v", EnvKillThisPid, os.Getpid()))
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

func tryKillProcess(pid string) error {
	processe, err := exec.Command("taskkill.exe", "/PID", pid, "/F").Output()
	if err != nil {
		return err
	}
	if strings.HasPrefix(string(processe), "SUCCESS") {
		return nil
	} else {
		return fmt.Errorf("failed to kill process")
	}
}
