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

// Build & publish image with ko
func KoPublish() error {
	// make sure you are logged in to the container registry
	imageTag, err := mage.GenerateImageTag()
	if err != nil {
		return fmt.Errorf("error generating image tag: %w", err)
	}
	err = mage.RunCmdWithLog(fmt.Sprintf("ko build --bare --tags %s ./cmd/aictl", imageTag))
	if err != nil {
		return fmt.Errorf("building image with ko: %w", err)
	}
	return nil
}
