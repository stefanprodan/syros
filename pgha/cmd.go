package main

import (
	"context"

	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func execId(timeout int) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "id", "-u", "-n").Output()
	if err != nil {
		log.Fatalf("id -u -n failed %s", err.Error())
	}
	return strings.Replace(string(output), "\n", "", -1)
}

func execPgStop(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "service", "postgresql", "stop").Run(); err != nil {
		log.Fatalf("service postgres stop failed %s", err.Error())
	}
}

func execRepmgrPromote(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "repmgr", "standby", "promote").Run(); err != nil {
		log.Warningf("repmgr standby promote failed %s", err.Error())
	} else {
		return
	}

	if err := exec.CommandContext(ctx, "repmgr", "standby", "promote").Run(); err != nil {
		log.Warningf("repmgr standby promote failed %s", err.Error())
		log.Info("trying repmgr standby switchover")
		execRepmgrSwitchover(timeout)
	}
}

func execRepmgrSwitchover(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "repmgr", "standby", "switchover").Run(); err != nil {
		log.Fatalf("repmgr standby switchover failed %s", err.Error())
	}
}

func execRepmgrVersion(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "repmgr", "-V").Output()
	if err != nil {
		log.Fatalf("repmgr -V failed %s", err.Error())
	}
	log.Info(strings.Replace(string(output), "\n", "", -1))
}
