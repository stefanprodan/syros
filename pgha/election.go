package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	consul "github.com/hashicorp/consul/api"
)

type Election struct {
	Key          string
	Session      string
	isLeader     bool
	consulClient *consul.Client
	consulLock   *consul.Lock
	stopChan     chan struct{}
	lockChan     chan struct{}
}

func NewElection(consulAddress string, key string, session string) (*Election, error) {
	cfg := consul.DefaultConfig()
	cfg.Address = consulAddress
	client, err := consul.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	lockOpt := &consul.LockOptions{
		Key: key,
		SessionOpts: &consul.SessionEntry{
			Name:      session,
			LockDelay: time.Duration(5 * time.Second),
			TTL:       "10s",
		},
	}
	lock, _ := client.LockOpts(lockOpt)

	e := &Election{
		Key:          key,
		Session:      session,
		isLeader:     false,
		consulClient: client,
		consulLock:   lock,
		stopChan:     make(chan struct{}, 1),
		lockChan:     make(chan struct{}, 1),
	}

	return e, nil
}

func (e *Election) Start() {
	stop := false
	for !stop {
		select {
		case <-e.stopChan:
			stop = true
		default:
			leader := e.GetLeader()
			if leader != "" {
				log.Infof("Entering follower state, leader is %s", leader)
			} else {
				log.Info("Entering candidate state, no leader found")
			}
			electionChan, err := e.consulLock.Lock(e.lockChan)
			if err != nil {
				log.Warnf("Failed to acquire election lock %s", err.Error())
			}
			if electionChan != nil {
				log.Info("Entering leader state")
				e.isLeader = true
				<-electionChan
				e.isLeader = false
				log.Warn("Leadership lost, releasing lock")
				e.consulLock.Unlock()
			} else {
				log.Info("Retrying election in 5s")
				time.Sleep(5000 * time.Millisecond)
			}
		}
	}
}

func (e *Election) GetLeader() string {
	kv, _, err := e.consulClient.KV().Get(e.Key, nil)
	if kv != nil && err == nil {
		sessionInfo, _, err := e.consulClient.Session().Info(kv.Session, nil)
		if err == nil {
			return sessionInfo.Name
		}
	}
	return ""
}

func (e *Election) Stop() {
	e.stopChan <- struct{}{}
	e.lockChan <- struct{}{}
	e.consulLock.Unlock()
}
