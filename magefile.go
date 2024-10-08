//go:build mage
// +build mage

package main

import (
	"fmt"
	"log/slog"
	"os"

	gbt "github.com/go-zen-chu/go-build-tools"
)

const currentVersion = "1.0.5"
const currentTagVersion = "v" + currentVersion

const imageRegistry = "amasuda"
const repository = "aictl"
const dockerFileLocation = "."

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

/*=======================
setup
=======================*/

// InstallDevTools installs required development tools for this project
func InstallDevTools() error {
	outMsg, errMsg, err := gbt.RunLongRunningCmdWithLog("go install go.uber.org/mock/mockgen@latest")
	if err != nil {
		return fmt.Errorf("installing mockgen: %w\nstdout: %s\nstderr: %s\n", err, outMsg, errMsg)
	}
	outMsg, errMsg, err = gbt.RunLongRunningCmdWithLog("go install github.com/goreleaser/goreleaser/v2@latest")
	if err != nil {
		return fmt.Errorf("installing goreleaser: %w\nstdout: %s\nstderr: %s\n", err, outMsg, errMsg)
	}
	return nil
}

/*=======================
workflow
=======================*/

func DockerLogin() error {
	return gbt.DockerLogin()
}

func DockerBuildLatest() error {
	return gbt.DockerBuildLatest(imageRegistry, repository, dockerFileLocation)
}

func DockerPublishLatest() error {
	return gbt.DockerPublishLatest(imageRegistry, repository)
}

func DockerBuildPublishWithGenTag() error {
	return gbt.DockerBuildPublishGeneratedImageTag(imageRegistry, repository, dockerFileLocation)
}

// GitPushTag pushes current tag to remote repository
func GitPushTag(releaseComment string) error {
	return gbt.GitPushTag(currentTagVersion, releaseComment)
}

const formulaTemplate = `class Aictl < Formula
    desc "Handy CLI tool to ask anything to generative AI in command line."
    homepage "https://github.com/go-zen-chu/aictl"
    version "%[1]s"
    
    on_macos do
        if Hardware::CPU.arm?
            url "https://github.com/go-zen-chu/aictl/releases/download/v%[1]s/aictl_Darwin_arm64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinArm64}}"
        else
            url "https://github.com/go-zen-chu/aictl/releases/download/v%[1]s/aictl_Darwin_x86_64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinX86_64}}"
        end
    end

    def install
        bin.install "aictl"
    end

    test do
        system "#{bin}/aictl", "--help"
    end
end
`

// UpdateFormula updates formula with current version for homebrew tap
func UpdateFormula() error {
	ft := fmt.Sprintf(formulaTemplate, currentVersion)
	return gbt.GenerateFormula(ft, "go-zen-chu", "homebrew-tools", "go-zen-chu", "aictl", currentTagVersion)
}
