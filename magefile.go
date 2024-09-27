//go:build mage
// +build mage

package main

import (
	"log/slog"
	"os"

	"github.com/go-zen-chu/aictl/internal/mage"
)

const imageRegistry = "amasuda"
const repository = "aictl"
const dockerFileLocation = "."

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

/*=======================
workflow
=======================*/

func DockerLogin() error {
	return mage.DockerLogin()
}

func DockerBuildLatest() error {
	return mage.DockerBuildLatest(imageRegistry, repository, dockerFileLocation)
}

func DockerPublishLatest() error {
	return mage.DockerPublishLatest(imageRegistry, repository, dockerFileLocation)
}

func DockerBuildPublishWithGenTag() error {
	return mage.DockerBuildPublishGeneratedImageTag(imageRegistry, repository, dockerFileLocation)
}
