package main

import (
	"fmt"
	"path"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
)

func createArtifactsDir(root string) (string, error) {
	session := sh.NewSession()
	session.SetDir(root)
	dir := fmt.Sprintf("deployctl/%d", time.Now().UnixNano())
	cmd := fmt.Sprintf("mkdir -p %s", dir)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "mkdir -p %s failed %s", dir, output)
	}
	return path.Join(root, dir), nil
}

func downloadArtifacts(url string, dir string) error {
	session := sh.NewSession()
	session.SetDir(dir)
	cmd := fmt.Sprintf("set -o pipefail; curl -sS %s | tar xvz", url)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "download %s failed %s", url, output)
	}
	return nil
}

func buildCompose(dir string, env string, host string) (string, error) {
	dir = path.Join(dir, "scripts")
	session := sh.NewSession()
	session.SetDir(dir)
	cmd := fmt.Sprintf("./build-host.sh -e %s -h %s -d .", env, host)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "./build-host.sh -e %s -h %s failed %s", env, host, output)
	}

	return path.Join(dir, host), nil
}
