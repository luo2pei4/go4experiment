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
		return nil, errors.New("invalid parameters, command is empty")
	}

	result = new(Result)
	cmd := exec.Command("/bin/bash", "-c", strings.Join(args, " "))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err = cmd.Start(); err != nil {
		return nil, err
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
	}

	if err != nil {
		return nil, err
	}

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
	result, err := ExecuteCmd(5*time.Second, "sleep", "10s")
	if err != nil {
		fmt.Println("execute error,", err.Error())
		return
	}

	if result.Code != 0 {
		fmt.Println("return code:", result.Code)
		fmt.Println(string(result.ErrOut))
		return
	}

	fmt.Println(string(result.StdOut))
}
