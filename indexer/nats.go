package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

func NewNatsConnection(servers string) (*nats.EncodedConn, error) {
	opts := nats.DefaultOptions
	opts.Url = servers
	host, _ := os.Hostname()
	opts.Name = fmt.Sprintf("syros-indexer-%s", host)
	opts.DisconnectedCB = func(nc *nats.Conn) {
		log.Warnf("Got disconnected from NATS %v", servers)
	}
	opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
	}
	opts.ClosedCB = func(nc *nats.Conn) {
		reason := "shutting down"
		if nc.LastError() != nil {
			reason = nc.LastError().Error()
		}
		log.Errorf("NATS connection closed. Reason: %v", reason)
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	enc, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return enc, err
}
