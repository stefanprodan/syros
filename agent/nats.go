package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

func NewNatsConnection(servers string, name string) (*nats.EncodedConn, error) {
	opts := nats.DefaultOptions
	opts.Url = servers
	opts.Name = name
	opts.DisconnectedCB = func(nc *nats.Conn) {
		log.Warnf("Got disconnected from NATS %v", servers)
	}
	opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
	}
	opts.ClosedCB = func(nc *nats.Conn) {
		if nc.LastError() != nil {
			log.Errorf("NATS connection closed. Reason: %v", nc.LastError().Error())
		}
	}
	opts.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, err error) {
		log.Error(err)
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

func natsPublish(natsCon string, subject string, v interface{}) error {
	nc, err := nats.Connect(natsCon)
	if err != nil {
		return err
	}
	enc, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	defer enc.Close()
	err = enc.Publish(subject, v)
	if err != nil {
		return err
	}

	return nil
}
