package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2/bson"
)

func (repo *Repository) DeploymentUpsert(dep models.Deployment) error {
	s := repo.Session.Copy()
	defer s.Close()

	// search for a release, update or insert
	r := s.DB(repo.Config.Database).C("releases")
	rel := models.Release{}
	rels := []models.Release{}
	err := r.Find(bson.M{"ticket_id": dep.TicketId}).All(&rels)
	if err != nil {
		log.Errorf("Repository DeploymentUpsert releases query failed %v", err)
		return err
	}

	if len(rels) < 1 {
		rel = models.Release{
			Id:       models.Hash(dep.TicketId),
			Begin:    time.Now().UTC(),
			End:      time.Now().UTC().Add(1 * time.Second),
			Name:     dep.TicketId,
			TicketId: dep.TicketId,
		}
	} else {
		rel = rels[0]
		rel.End = time.Now().UTC()
	}

	dlog := fmt.Sprintf("%v deployed on %v at %v env %v \n", dep.ServiceName, dep.HostName, time.Now().UTC(), dep.Environment)
	rel.Log += dlog
	rel.Deployments++

	_, err = r.UpsertId(rel.Id, &rel)
	if err != nil {
		log.Errorf("Repository DeploymentUpsert releases upsert failed %v", err)
		return err
	}

	dep.ReleaseId = rel.Id
	dep.Timestamp = time.Now().UTC()
	dep.Status = "Finished"
	dep.Id = models.Hash(fmt.Sprintf("%v%v%v", dep.TicketId, dep.ServiceName, dep.HostName))

	d := s.DB(repo.Config.Database).C("deployments")
	_, err = d.UpsertId(dep.Id, &dep)
	if err != nil {
		log.Errorf("Repository DeploymentUpsert deployments upsert failed %v", err)
		return err
	}

	return nil
}

func (repo *Repository) DeploymentStartUpsert(dep models.Deployment) error {
	s := repo.Session.Copy()
	defer s.Close()

	// search for a release, update or insert
	r := s.DB(repo.Config.Database).C("releases")
	rel := models.Release{}
	rels := []models.Release{}
	err := r.Find(bson.M{"ticket_id": dep.TicketId}).All(&rels)
	if err != nil {
		log.Errorf("Repository DeploymentStartUpsert releases query failed %v", err)
		return err
	}

	if len(rels) < 1 {
		rel = models.Release{
			Id:           models.Hash(dep.TicketId),
			Begin:        time.Now().UTC(),
			End:          time.Now().UTC().Add(1 * time.Second),
			Name:         dep.TicketId,
			TicketId:     dep.TicketId,
			Environments: dep.Environment,
		}
	} else {
		rel = rels[0]
		rel.End = time.Now().UTC()
		if !strings.Contains(rel.Environments, dep.Environment) {
			rel.Environments += fmt.Sprintf(", %v", dep.Environment)
		}
	}

	dlog := fmt.Sprintf("%v deploying on %v at %v env %v \n", dep.ServiceName, dep.HostName, time.Now().UTC(), dep.Environment)
	rel.Log += dlog

	_, err = r.UpsertId(rel.Id, &rel)
	if err != nil {
		log.Errorf("Repository DeploymentStartUpsert releases upsert failed %v", err)
		return err
	}

	return nil
}
