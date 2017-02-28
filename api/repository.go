package main

import (
	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/stefanprodan/syros/models"
)

type Repository struct {
	Config  *Config
	Session *r.Session
}

func NewRepository(config *Config) (*Repository, error) {

	session, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDB,
		Database: config.Database,
	})
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		Config:  config,
		Session: session,
	}
	return repo, nil
}

func (repo *Repository) AllHosts() ([]models.DockerHost, error) {
	cursor, err := r.Table("hosts").OrderBy(r.Asc("Collected"), r.OrderByOpts{Index: "Collected"}).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository AllHosts query failed %v", err)
		return nil, err
	}

	hosts := []models.DockerHost{}
	err = cursor.All(&hosts)
	if err != nil {
		log.Errorf("Repository AllHosts cursor failed %v", err)
		return nil, err
	}
	cursor.Close()

	return hosts, nil
}

func (repo *Repository) AllContainers() ([]models.DockerContainer, error) {
	cursor, err := r.Table("containers").OrderBy(r.Asc("Collected"), r.OrderByOpts{Index: "Collected"}).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository AllContainers query failed %v", err)
		return nil, err
	}

	containers := []models.DockerContainer{}
	err = cursor.All(&containers)
	if err != nil {
		log.Errorf("Repository AllContainers cursor failed %v", err)
		return nil, err
	}
	cursor.Close()

	return containers, nil
}
