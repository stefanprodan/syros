package main

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type PGStats struct {
	db       *sql.DB
	config   *Config
	stopChan chan struct{}
}

type ReplicationStats struct {
	Host      string    `json:"host"`
	Role      string    `json:"role"`
	Xlog      uint64    `json:"xlog"`
	Offset    uint64    `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

func NewPGStats(config *Config) (*PGStats, error) {
	db, err := sql.Open("postgres", config.PostgresURI)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	pg := &PGStats{
		db:       db,
		config:   config,
		stopChan: make(chan struct{}, 1),
	}
	return pg, nil
}

func (pg *PGStats) GetReplicationStats() (ReplicationStats, error) {
	stats := ReplicationStats{
		Host:      pg.config.Hostname,
		Timestamp: time.Now().UTC(),
	}
	var isInRecovery bool
	err := pg.db.QueryRow("SELECT pg_is_in_recovery()").Scan(&isInRecovery)
	if err != nil {
		return stats, err
	}

	var xlogLocation string
	if isInRecovery {
		stats.Role = "slave"
		err := pg.db.QueryRow("SELECT pg_last_xlog_receive_location();").Scan(&xlogLocation)
		if err != nil {
			return stats, err
		}
	} else {
		stats.Role = "master"
		err := pg.db.QueryRow("SELECT pg_current_xlog_location();").Scan(&xlogLocation)
		if err != nil {
			return stats, err
		}
	}

	xlog, err := strconv.ParseUint(strings.Split(xlogLocation, "/")[0], 16, 32)
	if err != nil {
		return stats, err
	}
	stats.Xlog = xlog

	offset, err := strconv.ParseUint(strings.Split(xlogLocation, "/")[1], 16, 32)
	if err != nil {
		return stats, err
	}
	stats.Offset = offset

	return stats, nil
}
