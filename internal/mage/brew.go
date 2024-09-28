package mage

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const tapRepo = "homebrew-tools"
const tapRepoUrl = "https://github.com/go-zen-chu/" + tapRepo + ".git"

type BrewFormula struct {
	ChecksumSHA256DarwinArm64  string
	ChecksumSHA256DarwinX86_64 string
}

func fileOrDirExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
}

func GenerateFormula(formulaTemplate, owner, repo, gitTag string) error {
	tmpl, err := template.New("formula").Parse(formulaTemplate)
	if err != nil {
		return fmt.Errorf("parse formula template: %w", err)
	}

	// clone tap repo to update formula
	if fileOrDirExists(tapRepo) {
		return fmt.Errorf("tap repo %s already exists. please remove before update formula", tapRepo)
	}
	outMsg, errMsg, err := RunLongRunningCmdWithLog(fmt.Sprintf("git clone %s", tapRepoUrl))
	if err != nil {
		return fmt.Errorf("git clone: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}

	var httpClient http.Client
	release, err := GetTagRelease(&httpClient, owner, repo, gitTag)
	if err != nil {
		return fmt.Errorf("get latest release: %w", err)
	}
	checksumMap, err := GetChecksumMap(&httpClient, release)
	if err != nil {
		return fmt.Errorf("get checksum map: %w", err)
	}
	bf := BrewFormula{}
	for filename, checksum := range checksumMap {
		if strings.Contains(filename, "Darwin") {
			if strings.Contains(filename, "arm64") {
				bf.ChecksumSHA256DarwinArm64 = checksum
			} else if strings.Contains(filename, "x86_64") {
				bf.ChecksumSHA256DarwinX86_64 = checksum
			}
		}
	}

	var bb bytes.Buffer
	err = tmpl.Execute(&bb, bf)
	if err != nil {
		return fmt.Errorf("execute formula template: %w", err)
	}
	formulaFilePath := fmt.Sprintf("%s/Formula/%s.rb", tapRepo, repo)
	err = os.WriteFile(formulaFilePath, bb.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("write formula file to %s: %w", formulaFilePath, err)
	}

	commitPushCmd := "bash -c " +
		"'" +
		"cd " + tapRepo +
		" && git config user.name \"GitHub Action\"" +
		" && git config user.email \"action@github.com\"" +
		" && git add --all" +
		" && git commit -m \"update formula to " + gitTag + "\"" +
		" && git push" +
		"'"
	outMsg, errMsg, err = RunLongRunningCmdWithLog(commitPushCmd)
	if err != nil {
		return fmt.Errorf("git commit and push: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}

	err = os.RemoveAll(tapRepo)
	if err != nil {
		return fmt.Errorf("remove tap repo %s: %w", tapRepo, err)
	}
	return nil
}
