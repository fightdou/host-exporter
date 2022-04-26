package pkg

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"
)

func Execute(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	stdoutStderr, err := cmd.CombinedOutput()
	return stdoutStderr, err
}

func RunCommandWithTimeout(timeout time.Duration, command string, args ...string) (stdout []byte, isKilled bool) {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Start()
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(timeout)
	select {
	case <-after:
		cmd.Process.Signal(syscall.SIGINT)
		time.Sleep(10 * time.Millisecond)
		cmd.Process.Kill()
		isKilled = true
	case <-done:
		isKilled = false
	}
	stdout = bytes.TrimSpace(stdoutBuf.Bytes())
	return
}
