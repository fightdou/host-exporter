package pkg

import "os/exec"

func Execute(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	stdoutStderr, err := cmd.CombinedOutput()
	return stdoutStderr, err
}
