package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

func NewNatsConnection(servers string) (*nats.Conn, error) {
	nc, err := nats.Connect(servers,
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Warnf("Got disconnected from NATS %v", servers)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Errorf("NATS connection closed. Reason: %q", nc.LastError())
		}),
	)
	return nc, err
}
