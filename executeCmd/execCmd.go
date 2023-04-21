package main

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Result struct {
	ReturnCode  int
	StdOutBuf   []byte
	ErrOutBuf   []byte
	CommandLine string
}

func ExecuteCmd(dur time.Duration, args ...string) (result *Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if len(args) == 0 {
		return nil, errors.New("invalid parameters, command is empty")
	}

	result = new(Result)
	cmd := exec.Command(strings.Join(args, " "))

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
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	if err != nil {
		return nil, err
	}

	result.ReturnCode = cmd.ProcessState.ExitCode()
	if stdout.Len() > 0 {
		result.StdOutBuf = stdout.Bytes()
	}
	if stderr.Len() > 0 {
		result.ErrOutBuf = stderr.Bytes()
	}

	return
}
