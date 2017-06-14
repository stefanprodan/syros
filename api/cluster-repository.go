package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2/bson"
)

func (repo *Repository) AllClusterHealthChecks() ([]models.ClusterHealthCheck, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("cluster_checks")
	checks := []models.ClusterHealthCheck{}
	err := c.Find(nil).Sort("-collected").All(&checks)
	if err != nil {
		log.Errorf("Repository AllClusterHealthChecks query failed %v", err)
		return nil, err
	}

	return checks, nil
}

func (repo *Repository) ClusterHealthCheckLog(checkId string) ([]models.ClusterHealthCheckLog, []models.HealthCheckStats, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("cluster_checks_log")
	logs := []models.ClusterHealthCheckLog{}
	err := c.Find(bson.M{"check_id": checkId}).Sort("-begin").Limit(500).All(&logs)
	if err != nil {
		log.Errorf("Repository ClusterHealthCheckLog checks_log query failed %v", err)
		return nil, nil, err
	}

	k := s.DB(repo.Config.Database).C("cluster_checks")
	current := models.ClusterHealthCheck{}
	err = k.FindId(checkId).One(&current)
	if err != nil {
		log.Errorf("Repository ClusterHealthCheckLog checks query failed %v", err)
		return nil, nil, err
	}

	// add current status to logs
	cur := models.NewClusterHealthCheckLog(current, current.Since, time.Now().UTC())
	logs = append(logs, cur)

	last30d := time.Now().UTC().Add((-30 * 24) * time.Hour)
	stats := []models.HealthCheckStats{}

	pipeline := []bson.M{
		{"$match": bson.M{
			"check_id": checkId,
			"begin":    bson.M{"$gt": last30d},
		}},
		{"$group": bson.M{
			"_id":      "$status",
			"count":    bson.M{"$sum": 1},
			"duration": bson.M{"$sum": "$duration"},
		}},
	}

	pipe := c.Pipe(pipeline)
	err = pipe.All(&stats)
	if err != nil {
		log.Errorf("Repository ClusterHealthCheckLog pipeline failed %v", err)
		return nil, nil, err
	}

	// add current status to stats
	found := false
	for i, stat := range stats {
		if stat.Status == cur.Status {
			stats[i].Count++
			stats[i].Duration += cur.Duration
			found = true
		}
	}
	if !found {
		stat := models.HealthCheckStats{
			Status:   cur.Status,
			Count:    1,
			Duration: cur.Duration,
		}
		stats = append(stats, stat)
	}

	return logs, stats, nil
}
