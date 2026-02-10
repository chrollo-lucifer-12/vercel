package main

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"runtime"

	"github.com/chrollo-lucifer-12/build-server/logs"
	"github.com/google/uuid"
)

func RunNpmCommand(
	ctx context.Context,
	dir string,
	dispatcher *logs.LogDispatcher,
	deploymentID uuid.UUID,
	args ...string,
) error {
	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	cmd := exec.CommandContext(ctx, npm, args...)
	cmd.Dir = dir

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	stream := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			dispatcher.Push(deploymentID, scanner.Text())
		}
	}

	go stream(stdoutPipe)
	go stream(stderrPipe)

	return cmd.Wait()
}
