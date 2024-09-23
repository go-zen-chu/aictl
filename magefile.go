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
	imageTag, err := mage.GenerateImageTag()
	if err != nil {
		return fmt.Errorf("error generating image tag: %w", err)
	}
	// make sure you are logged in to the container registry
	// TIPS: ko build would add `-<md5>`` to image name without --base-import-paths flag
	out, errMsg, err := mage.RunLongRunningCmd(fmt.Sprintf("ko build --base-import-paths --tags %s ./cmd/aictl", imageTag))
	if err != nil {
		return fmt.Errorf("building image with ko: %w, stdout log: %s, stderr log: %s", err, out, errMsg)
	}
	return nil
}

// Build & publish latest tag with ko
func KoPublishLatest() error {
	out, errMsg, err := mage.RunLongRunningCmd("ko build --base-import-paths --tags latest ./cmd/aictl")
	if err != nil {
		return fmt.Errorf("building image with ko: %w, stdout log: %s, stderr log: %s", err, out, errMsg)
	}
	return nil
}
