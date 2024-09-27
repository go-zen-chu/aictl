//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/go-zen-chu/aictl/internal/mage"
)

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

/*=======================
setup
=======================*/

// Install tools for developping this repository
func InstallDevTools() error {
	// make sure golang is already installed
	return mage.RunCmdWithLog("go install github.com/google/ko@latest")
}

/*=======================
workflow
=======================*/

// Push to docker.io
func DockerPublish() error {
	outMsg, errMsg, err := mage.RunLongRunningCmd("docker push amasuda/aictl:latest")
	if err != nil {
		return fmt.Errorf("pushing to docker: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}
	return nil
}
