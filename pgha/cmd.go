package main

import (
	"context"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
)

func execPgStop(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "service", "postgres", "stop").Run(); err != nil {
		log.Fatalf("service postgres stop failed %s", err.Error())
	}
}

func execRepmgrPromote(timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "repmgr", "standby", "promote").Run(); err != nil {
		log.Fatalf("repmgr standby promote failed %s", err.Error())
	}
}
