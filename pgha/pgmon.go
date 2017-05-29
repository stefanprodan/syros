package main

import (
	"database/sql"

	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

type PGMonitor struct {
	db       *sql.DB
	status   *Status
	stopChan chan struct{}
}

func NewPGMonitor(uri string, status *Status) (*PGMonitor, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	pg := &PGMonitor{
		db:       db,
		status:   status,
		stopChan: make(chan struct{}, 1),
	}
	return pg, nil
}

func (pg *PGMonitor) IsMaster() (bool, error) {
	var isInRecovery bool
	err := pg.db.QueryRow("SELECT pg_is_in_recovery()").Scan(&isInRecovery)
	if err != nil {
		return false, err
	}

	return !isInRecovery, nil
}

func (pg *PGMonitor) Start() {
	running := true
	for running {
		select {
		case <-pg.stopChan:
			running = false
		default:
			isMaster, err := pg.IsMaster()
			if err != nil {
				log.Warnf("Failed to acquire PG state %s", err.Error())
				pg.status.SetPostgresStatus(false)
			} else {
				pg.status.SetPostgresStatus(isMaster)
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func (pg *PGMonitor) Stop() {
	pg.stopChan <- struct{}{}
}
