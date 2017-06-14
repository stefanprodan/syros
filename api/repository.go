package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2"
)

type Repository struct {
	Config  *Config
	Session *mgo.Session
}

func NewRepository(config *Config) (*Repository, error) {
	cluster := strings.Split(config.MongoDB, ",")
	dialInfo := &mgo.DialInfo{
		Addrs:    cluster,
		Database: config.Database,
		Timeout:  10 * time.Second,
		FailFast: true,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	repo := &Repository{
		Config:  config,
		Session: session,
	}

	return repo, nil
}

func (repo *Repository) AllSyrosServices() ([]models.SyrosService, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("syros_services")
	services := []models.SyrosService{}
	err := c.Find(nil).Sort("-collected").All(&services)
	if err != nil {
		log.Errorf("Repository AllSyrosServices query failed %v", err)
		return nil, err
	}

	return services, nil
}
