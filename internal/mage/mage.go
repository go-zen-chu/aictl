package mage

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// split those are (not space | start and end with double quote | start and end with single quote)
var cmdSplitRegex = regexp.MustCompile(`[^\s"']+|"[^"]*"|'[^']*'`)

func splitCmd(cmd string) []string {
	matches := cmdSplitRegex.FindAllString(cmd, -1)
	if matches == nil {
		slog.Error("could not split command", "cmd", cmd)
		panic("could not split command with white space")
	}
	for i, m := range matches {
		matches[i] = strings.Trim(m, `"'`)
	}
	return matches
}

func runCmd(cmd string, preRunCmd func(cmd string)) (string, error) {
	if cmd == "" {
		return "", fmt.Errorf("command string is empty")
	}
	cmdSplit := splitCmd(cmd)
	c := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	if preRunCmd != nil {
		preRunCmd(cmd)
	}
	out, err := c.CombinedOutput()
	return string(out), err
}

// RunCmdWithResult runs a command and returns the result and error of the command
func RunCmdWithResult(cmd string) (string, error) {
	return runCmd(cmd, nil)
}

// RunCmdWithResultWithLog runs a command and returns the result and error of the command. It logs which command was run before running it
// **Do not use this function with commands that have sensitive information**
func RunCmdWithResultWithLog(cmd string) (string, error) {
	return runCmd(cmd, func(cmd string) {
		slog.Info("Running command", "cmd", cmd)
	})
}

// RunCmdWithLog runs a command and logs the result
// **Do not use this function with commands that have sensitive information**
func RunCmdWithLog(cmd string) error {
	out, err := RunCmdWithResultWithLog(cmd)
	if err != nil {
		return fmt.Errorf("error running command: %w\nerror log: %s", err, out)
	}
	if out == "" {
		out = "[command result was empty string]"
	}
	slog.Info(out)
	return nil
}

func runLongRunningCmd(cmd string, preRunCmd func(cmd string)) (string, string, error) {
	cmdSplit := splitCmd(cmd)
	c := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	cmdStdoutBuffer := bytes.NewBufferString("")
	cmdStderrBuffer := bytes.NewBufferString("")
	// write to both stdout and string buffer
	cmdStdoutMultiWriter := io.MultiWriter(os.Stdout, cmdStdoutBuffer)
	// write to both stderr and string buffer
	cmdStderrMultiWriter := io.MultiWriter(os.Stderr, cmdStderrBuffer)
	c.Stdout = cmdStdoutMultiWriter
	c.Stderr = cmdStderrMultiWriter
	if preRunCmd != nil {
		preRunCmd(cmd)
	}
	err := c.Run()
	if err != nil {
		return cmdStdoutBuffer.String(), cmdStderrBuffer.String(), fmt.Errorf("run command: %w", err)
	}
	// return both results with string
	return cmdStdoutBuffer.String(), cmdStderrBuffer.String(), nil
}

// RunLongRunningCmd runs a command that takes long time and returns the result and error of the command
func RunLongRunningCmd(cmd string) (string, string, error) {
	return runLongRunningCmd(cmd, nil)
}

// RunLongRunningCmdWithLog runs a command that takes long time and returns the result and error of the command with logs
// **Do not use this function with commands that have sensitive information**
func RunLongRunningCmdWithLog(cmd string) (string, string, error) {
	return runLongRunningCmd(cmd, func(cmd string) {
		slog.Info("Running long running command", "cmd", cmd)
	})
}

func getGitCommitShortHash() (string, error) {
	out, err := RunCmdWithResult("git rev-parse --short HEAD")
	if err != nil {
		return "", fmt.Errorf("getting git commit short hash: %w\nerror log: %s", err, out)
	}
	commitShortHash := strings.TrimSuffix(out, "\n")
	return commitShortHash, err
}

func getCurrentBranch() (string, error) {
	out, err := RunCmdWithResult("git symbolic-ref --short HEAD")
	if err != nil {
		return "", fmt.Errorf("getting current branch: %w\nerror log: %s", err, out)
	}
	branch := strings.TrimSuffix(out, "\n")
	return branch, err
}

func getCurrentDateTimeUTC() string {
	return time.Now().UTC().Format("2006-01-02T15-04-05Z")
}

// GenerateImageTag generates an image tag with following format:
// {branch}_{commit short hash}_{current datetime in UTC}
func GenerateImageTag() (string, error) {
	commitHash, err := getGitCommitShortHash()
	if err != nil {
		return "", fmt.Errorf("generating image tag: %w", err)
	}
	branch, err := getCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("generating image tag: %w", err)
	}
	currentDateTime := getCurrentDateTimeUTC()
	return strings.Join([]string{branch, commitHash, currentDateTime}, "_"), nil
}
