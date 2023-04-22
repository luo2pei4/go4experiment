package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Result struct {
	Code   int
	StdOut []byte
	ErrOut []byte
}

func ExecuteCmd(dur time.Duration, args ...string) (result *Result, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if len(args) == 0 {
		err = errors.New("invalid parameters, command is empty")
		return
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", strings.Join(args, " "))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err = cmd.Start(); err != nil {
		return
	}

	notify := make(chan struct{})
	go func() {
		err = cmd.Wait()
		close(notify)
	}()

	select {
	case <-notify:
	case <-ctx.Done():
		if err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err == nil {
			err = errors.New("execution timeout")
		}
		// return directly when execution timeout
		return
	}

	result = new(Result)
	result.Code = cmd.ProcessState.ExitCode()
	if stdout.Len() > 0 {
		result.StdOut = stdout.Bytes()
	}
	if stderr.Len() > 0 {
		result.ErrOut = stderr.Bytes()
	}

	return
}

func main() {
	result, err := ExecuteCmd(5*time.Second, "sleep", "8s")
	if err != nil {
		fmt.Println("execute error,", err.Error())
	}

	if result != nil {
		fmt.Println("exit code:", result.Code)
		if len(result.ErrOut) > 0 {
			fmt.Println("cmd errout:", string(result.ErrOut))
		}
		if len(result.StdOut) > 0 {
			fmt.Println("cmd stdout", string(result.StdOut))
		}
	}
}
