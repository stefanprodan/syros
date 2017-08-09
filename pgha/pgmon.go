package main

import (
	"database/sql"

	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type PGMonitor struct {
	db       *sql.DB
	status   *Status
	stopChan chan struct{}
}

func NewPGMonitor(uri string, status *Status) (*PGMonitor, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, errors.Wrap(err, "Postgres init failed")
	}
	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Postgres ping failed")
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
		return false, errors.Wrap(err, "Query pg_is_in_recovery failed")
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
			isMaster, err := pg.GetMasterWithRetry(5, 1)
			if err != nil {
				log.Fatalf("Failed to acquire PG state %s", err.Error())
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

func (pg *PGMonitor) GetMasterWithRetry(retry int, wait int) (bool, error) {
	var leader bool
	var err error
	for retry > 0 {
		leader, err = pg.IsMaster()
		if err != nil {
			retry--
			log.Warnf("Failed to acquire PG state retrying %v after %v seconds", retry, wait)
			time.Sleep(time.Duration(wait) * time.Second)
		} else {
			return leader, nil
		}
	}

	return leader, err
}
