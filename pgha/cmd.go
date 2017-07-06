package main

import (
	"context"

	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

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
		log.Fatalf("repmgr standby promote failed %s", err.Error())
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
