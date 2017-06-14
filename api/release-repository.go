package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2/bson"
)

func (repo *Repository) AllReleases() ([]models.Release, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("releases")
	rels := []models.Release{}
	err := c.Find(nil).Sort("end").Limit(1000).All(&rels)
	if err != nil {
		log.Errorf("Repository AllReleases query failed %v", err)
		return nil, err
	}

	return rels, nil
}

func (repo *Repository) ReleaseDeployments(releaseId string) ([]models.Deployment, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("deployments")
	deployments := []models.Deployment{}
	err := c.Find(bson.M{"release_id": releaseId}).Sort("-end").All(&deployments)
	if err != nil {
		log.Errorf("Repository ReleaseDeployments query failed %v", err)
		return nil, err
	}

	return deployments, nil
}
