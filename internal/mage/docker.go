package mage

import (
	"fmt"
	"os"
)

// DockerLogin logs in to docker.io
func DockerLogin() error {
	user := os.Getenv("DOCKER_USERNAME")
	if user == "" {
		return fmt.Errorf("DOCKER_USERNAME is not set")
	}
	pswd := os.Getenv("DOCKER_PASSWORD")
	if pswd == "" {
		return fmt.Errorf("DOCKER_PASSWORD is not set")
	}
	outMsg, errMsg, err := RunLongRunningCmd(fmt.Sprintf("docker login -u %s -p %s", user, pswd))
	if err != nil {
		return fmt.Errorf("docker login failed: %w\nstdout: %s\nstderr: %s", err, outMsg, errMsg)
	}
	return nil
}

// DockerBuild builds docker image
func DockerBuild(registry string, repository string, tag string, dockerFileLocation string) error {
	buildCmd := fmt.Sprintf("docker build -t %s/%s:%s %s", registry, repository, tag, dockerFileLocation)
	outMsg, errMsg, err := RunLongRunningCmdWithLog(buildCmd)
	if err != nil {
		return fmt.Errorf("building docker image (%s): %w\nstdout: %s\nstderr: %s", buildCmd, err, outMsg, errMsg)
	}
	return nil
}

// DockerBuildLatest builds docker image with latest tag
func DockerBuildLatest(registry string, repository string, dockerFileLocation string) error {
	return DockerBuild(registry, repository, "latest", dockerFileLocation)
}

// DockerPublish pushes image
func DockerPublish(registry string, repository string, tag string) error {
	pushCmd := fmt.Sprintf("docker push %s/%s:%s", registry, repository, tag)
	outMsg, errMsg, err := RunLongRunningCmdWithLog(pushCmd)
	if err != nil {
		return fmt.Errorf("pushing to docker (%s): %w\nstdout: %s\nstderr: %s", pushCmd, err, outMsg, errMsg)
	}
	return nil
}

// DockerPublishLatest pushes image with latest tag
func DockerPublishLatest(registry string, repository string) error {
	return DockerPublish(registry, repository, "latest")
}

// DockerBuildPublishLatest builds and pushes image with latest tag
func DockerBuildPublishLatest(registry string, repository string, dockerFileLocation string) error {
	if err := DockerBuildLatest(registry, repository, dockerFileLocation); err != nil {
		return fmt.Errorf("building latest docker image: %w", err)
	}
	if err := DockerPublishLatest(registry, repository); err != nil {
		return fmt.Errorf("pushing latest docker image: %w", err)
	}
	return nil
}

// DockerBuildPublishGeneratedImageTag builds and pushes image with generated tag
func DockerBuildPublishGeneratedImageTag(registry string, repository string, dockerFileLocation string) error {
	tag, err := GenerateImageTag()
	if err != nil {
		return fmt.Errorf("generating image tag: %w", err)
	}
	if err := DockerBuild(registry, repository, tag, dockerFileLocation); err != nil {
		return fmt.Errorf("building docker image with generated tag: %w", err)
	}
	if err := DockerPublish(registry, repository, tag); err != nil {
		return fmt.Errorf("pushing docker image with generated tag: %w", err)
	}
	return nil
}
