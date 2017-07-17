package main

import (
	"fmt"
	"log"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
)

type PostgresDeploy struct {
	Env      string
	HostTo   string
	Service  string
	Dir      string
	Ssh      *SshClient
	Url      string
	User     string
	Password string
	Location string
	Database string
}

func (m PostgresDeploy) Migrate() error {
	log.Printf("Starting %s migration to %s", m.Service, m.HostTo)

	session := sh.NewSession()
	session.SetDir(m.Dir)
	cmd := fmt.Sprintf("./flyway -url=%s -user=%s -password=%s -locations=filesystem:%s -outOfOrder=true migrate -X",
		m.Url, m.User, m.Password, m.Location)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "flyway failed %s", string(output))
	}
	log.Print(string(output))

	return nil
}
