package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
)

func containerExists(dockerHost string, image string) (bool, error) {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("set -o pipefail; docker ps --format '{{.Names}}' -a| awk '$1 ~ /%s$/ {print $1}'", image)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return false, errors.Wrapf(err, "container lookup on %s for %s failed %s", dockerHost, image, output)
	}
	if len(output) < 1 {
		return false, nil
	}

	return true, nil
}

func containerIsRunning(dockerHost string, image string) (bool, error) {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("set -o pipefail; docker ps --format '{{.Names}}'| awk '$1 ~ /%s$/ {print $1}'", image)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return false, errors.Wrapf(err, "container lookup on %s for %s failed %s", dockerHost, image, output)
	}
	if len(output) < 1 {
		return false, nil
	}

	return true, nil
}

func containerRename(dockerHost string, name string, newName string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("docker rename %s %s", name, newName)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker rename %s %s failed %s", name, newName, output)
	}
	return nil
}

func containerStop(dockerHost string, name string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("docker stop %s", name)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker stop %s failed %s", name, output)
	}
	return nil
}

func containerStart(dockerHost string, name string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("docker start %s", name)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker start %s failed %s", name, output)
	}
	return nil
}

func containerRemove(dockerHost string, name string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("docker rm -f %s", name)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker rm -f %s failed %s", name, output)
	}
	return nil
}

func imageGetTag(dockerHost string, image string) (string, error) {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	cmd := fmt.Sprintf("set -o pipefail; docker ps --format '{{.Image}}' -a| awk '$1 ~ /%s/ {print $1}'| awk -F: '{print $NF}'| head -n1", image)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "acquiring image tag from %s for %s failed %s", dockerHost, image, output)
	}
	if len(output) < 1 {
		return "", errors.Errorf("no image tag found on %s for %s", dockerHost, image)
	}
	tag := strings.TrimSpace(fmt.Sprintf("%s", output))
	return tag, nil
}

func imagePurge(dockerHost string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	output, err := session.Command("/bin/sh", "-c", "docker image prune -af").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker image prune -af %s", output)
	}
	return nil
}

func imagePull(dockerHost string, dir string, image string, tag string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	session.SetDir(dir)
	ver := "VER_" + strings.ToUpper(strings.Replace(image, "-", "_", -1))
	session.SetEnv(ver, tag)
	cmd := fmt.Sprintf("docker-compose pull %s", image)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker-compose pull %s failed %s", image, output)
	}
	return nil
}

func containerCreate(dockerHost string, dir string, image string, tag string, ticket string) error {
	session := sh.NewSession()
	session.SetEnv("DOCKER_HOST", dockerHost)
	session.SetEnv("TICKET", ticket)
	ver := "VER_" + strings.ToUpper(strings.Replace(image, "-", "_", -1))
	session.SetEnv(ver, tag)
	session.SetDir(dir)
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	cmd := fmt.Sprintf("docker-compose -p %s create --force-recreate %s", id, image)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker-compose -p %s create --force-recreate %s failed %s", id, image, output)
	}
	return nil
}
