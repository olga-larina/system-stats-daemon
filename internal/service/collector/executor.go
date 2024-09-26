package collector

import "os/exec"

type CommandExecutor interface {
	Execute(name string, arg ...string) ([]byte, error)
}

type OsCommandExecutor struct{}

func NewOsCommandExecutor() *OsCommandExecutor {
	return &OsCommandExecutor{}
}

func (r *OsCommandExecutor) Execute(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}
