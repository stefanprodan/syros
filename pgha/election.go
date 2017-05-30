package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	consul "github.com/hashicorp/consul/api"
)

type Election struct {
	key          string
	session      string
	ttl          string
	status       *Status
	consulClient *consul.Client
	consulLock   *consul.Lock
	stopChan     chan struct{}
	lockChan     chan struct{}
}

func NewElection(consulAddress string, ttl string, prefix string, session string, status *Status) (*Election, error) {
	key := prefix + "/leader/election"
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
			TTL:       ttl,
		},
	}
	lock, _ := client.LockOpts(lockOpt)

	e := &Election{
		key:          key,
		session:      session,
		status:       status,
		consulClient: client,
		consulLock:   lock,
		stopChan:     make(chan struct{}, 1),
		lockChan:     make(chan struct{}, 1),
	}

	return e, nil
}

func (e *Election) Start() {
	runElection := true
	for runElection {
		select {
		case <-e.stopChan:
			runElection = false
		default:
			leader, err := e.GetLeaderWithRetry(5, 1)
			if err != nil {
				//TODO: stop pg service if current node is master
				log.Warnf("Consul is unreachable %s", err.Error())
				e.status.SetConsulStatus(false, FaultedCode, err.Error())
			}
			if leader != "" {
				log.Infof("Entering follower state, leader is %s", leader)
				e.status.SetConsulStatus(false, FollowerCode, fmt.Sprintf("follower of %s", leader))
			} else {
				log.Info("Entering candidate state, no leader found")
			}
			electionChan, err := e.consulLock.Lock(e.lockChan)
			if err != nil {
				log.Warnf("Failed to acquire election lock %s", err.Error())
				e.status.SetConsulStatus(false, FaultedCode, err.Error())
			}
			if electionChan != nil {
				log.Info("Entering leader state")
				e.status.SetConsulStatus(true, LeaderCode, "leader")
				<-electionChan
				//TODO: switch pg to slave mode or keep retrying, detect shutdown mode, check for new leader
				log.Warn("Leadership lost, releasing lock")
				e.status.SetConsulStatus(false, FaultedCode, "leadership lost")
				e.consulLock.Unlock()
			} else {
				log.Info("Retrying election in 5s")
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (e *Election) GetLeader() (string, error) {
	kv, _, err := e.consulClient.KV().Get(e.key, nil)
	if err != nil {
		return "", err
	}

	if kv != nil {
		sessionInfo, _, err := e.consulClient.Session().Info(kv.Session, nil)
		if sessionInfo != nil && err == nil {
			return sessionInfo.Name, nil
		}
	}
	return "", nil
}

func (e *Election) GetLeaderWithRetry(retry int, wait int) (string, error) {
	var leader string
	var err error
	for retry > 0 {
		leader, err = e.GetLeader()
		if err != nil {
			retry--
			log.Warnf("Consul is unreachable retrying %v after %v seconds", retry, wait)
			time.Sleep(time.Duration(wait) * time.Second)
		} else {
			return leader, nil
		}
	}

	return leader, err
}

func (e *Election) Fallback() error {
	e.lockChan <- struct{}{}
	return e.consulLock.Unlock()
}

func (e *Election) Stop() {
	e.stopChan <- struct{}{}
	e.lockChan <- struct{}{}
	e.consulLock.Unlock()
}
