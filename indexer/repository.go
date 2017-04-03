package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
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

func (repo *Repository) Initialize() {
	repo.CreateIndex("hosts", "environment")
	repo.CreateIndex("hosts", "collected")
	repo.CreateIndex("containers", "host_id")
	repo.CreateIndex("containers", "environment")
	repo.CreateIndex("containers", "collected")
	repo.CreateIndex("checks", "host_id")
	repo.CreateIndex("checks", "environment")
	repo.CreateIndex("checks", "collected")
	repo.CreateIndex("syros_services", "environment")
	repo.CreateIndex("syros_services", "collected")
}

func (repo *Repository) CreateIndex(col string, index string) {
	c := repo.Session.DB(repo.Config.Database).C(col)
	err := c.EnsureIndexKey(index)

	if err != nil {
		log.Fatalf("MongoDB index %v init failed %v", index, err)
	}
}

func (repo *Repository) HostUpsert(host models.DockerHost) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("hosts")

	_, err := c.UpsertId(host.Id, &host)
	if err != nil {
		log.Errorf("Repository hosts upsert failed %v", err)
	}
}

func (repo *Repository) ContainerUpsert(container models.DockerContainer) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("containers")

	_, err := c.UpsertId(container.Id, &container)
	if err != nil {
		log.Errorf("Repository containers upsert failed %v", err)
	}
}

func (repo *Repository) ContainersUpsert(containers []models.DockerContainer) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("containers")

	for _, container := range containers {
		_, err := c.UpsertId(container.Id, &container)
		if err != nil {
			log.Errorf("Repository containers upsert failed %v", err)
		}
	}
}

func (repo *Repository) ChecksUpsert(checks []models.ConsulHealthCheck) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("checks")

	for _, check := range checks {
		res := models.ConsulHealthCheck{}
		err := c.FindId(check.Id).One(&res)
		if err != nil {
			if err.Error() == "not found" {
				check.Since = check.Collected
				_, err = c.UpsertId(check.Id, &check)
				if err != nil {
					log.Errorf("Repository checks insert failed %v", err)
				}
				return
			} else {
				log.Errorf("Repository checks find by id failed %v", err)
			}
		}
		if res.Status != check.Status {
			check.Since = check.Collected
		} else {
			check.Since = res.Since
		}
		_, err = c.UpsertId(check.Id, &check)
		if err != nil {
			log.Errorf("Repository checks upsert failed %v", err)
		}
	}
}

func (repo *Repository) SyrosServiceUpsert(service models.SyrosService) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("syros_services")

	_, err := c.UpsertId(service.Id, &service)
	if err != nil {
		log.Errorf("Repository syros_services upsert failed %v", err)
	}
}

// Removes stale records
func (repo *Repository) RunGarbageCollector(cols []string) {
	if repo.Config.DatabaseStale > 0 {
		log.Infof("Stating repository GC interval %v minutes", repo.Config.DatabaseStale)
		go func(stale int) {

			for true {
				s := repo.Session.Copy()
				for _, col := range cols {
					c := s.DB(repo.Config.Database).C(col)
					info, err := c.RemoveAll(
						bson.M{
							"collected": bson.M{
								"$lt": time.Now().Add(-time.Duration(stale) * time.Minute).UTC(),
							},
						})
					if err != nil {
						log.Errorf("Repository GC for col %v query failed %v", col, err)
					} else {
						if info.Removed > 0 {
							log.Infof("Repository GC removed %v from %v", info.Removed, col)
						}
					}
				}
				s.Close()
				time.Sleep(60 * time.Second)
			}

		}(repo.Config.DatabaseStale)
	}
}
