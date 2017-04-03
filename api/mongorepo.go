package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type MongoRepository struct {
	Config  *Config
	Session *mgo.Session
}

func NewMongoRepository(config *Config) (*MongoRepository, error) {
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

	repo := &MongoRepository{
		Config:  config,
		Session: session,
	}

	return repo, nil
}

func (repo *MongoRepository) AllEnvironments() ([]string, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("hosts")
	var result []string
	err := c.Find(nil).Distinct("environment", &result)
	if err != nil {
		log.Errorf("Repository AllEnvironments query failed %v", err)
		return nil, err
	}

	return result, nil
}

type HostStats struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
}

func (repo *MongoRepository) EnvironmentHostContainerSum() ([]models.EnvironmentStats, error) {
	s := repo.Session.Copy()
	defer s.Close()

	h := s.DB(repo.Config.Database).C("hosts")
	var all []string
	err := h.Find(nil).Distinct("environment", &all)
	if err != nil {
		log.Errorf("Repository EnvironmentHostContainerSum query failed %v", err)
		return nil, err
	}

	environments := []models.EnvironmentStats{}

	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":                "$environment",
			"hosts":              bson.M{"$sum": 1},
			"containers_running": bson.M{"$sum": "$containers_running"},
			"ncpu":               bson.M{"$sum": "$ncpu"},
			"mem_total":          bson.M{"$sum": "$mem_total"},
		}},
	}

	pipe := h.Pipe(pipeline)
	err = pipe.All(&environments)
	if err != nil {
		log.Errorf("Repository EnvironmentHostContainerSum pipeline failed %v", err)
		return nil, err
	}

	return environments, nil
}

func (repo *MongoRepository) AllHosts() ([]models.DockerHost, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("hosts")
	hosts := []models.DockerHost{}
	err := c.Find(nil).Sort("-collected").All(&hosts)
	if err != nil {
		log.Errorf("Repository AllHosts cursor failed %v", err)
		return nil, err
	}

	return hosts, nil
}

func (repo *MongoRepository) HostContainers(hostID string) (*models.DockerPayload, error) {
	s := repo.Session.Copy()
	defer s.Close()

	h := s.DB(repo.Config.Database).C("hosts")
	host := models.DockerHost{}
	err := h.FindId(hostID).One(&host)
	if err != nil {
		log.Errorf("Repository HostContainers query failed for hostID %v %v", hostID, err)
		return nil, err
	}

	c := s.DB(repo.Config.Database).C("containers")
	containers := []models.DockerContainer{}
	err = c.Find(bson.M{"host_id": hostID}).Sort("-collected").All(&containers)
	if err != nil {
		log.Errorf("Repository HostContainers query containers All for hostID %v failed %v", hostID, err)
		return nil, err
	}

	payload := &models.DockerPayload{
		Host:       host,
		Containers: containers,
	}

	return payload, nil
}

func (repo *MongoRepository) EnvironmentContainers(env string) (*models.DockerPayload, error) {
	s := repo.Session.Copy()
	defer s.Close()

	h := s.DB(repo.Config.Database).C("hosts")
	hosts := []models.DockerHost{}
	err := h.Find(bson.M{"environment": env}).All(&hosts)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers query hosts for env %v failed %v", env, err)
		return nil, err
	}

	envStats := models.DockerHost{}

	for _, host := range hosts {
		envStats.ContainersRunning += host.ContainersRunning
		envStats.Containers++
		envStats.NCPU += host.NCPU
		envStats.MemTotal += host.MemTotal
	}

	c := s.DB(repo.Config.Database).C("containers")
	containers := []models.DockerContainer{}
	err = c.Find(bson.M{"environment": env}).Sort("-collected").All(&containers)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers query containers All for env %v failed %v", env, err)
		return nil, err
	}

	payload := &models.DockerPayload{
		Host:       envStats,
		Containers: containers,
	}

	return payload, nil
}

func (repo *MongoRepository) Container(containerID string) (*models.DockerPayload, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("containers")
	container := models.DockerContainer{}
	err := c.FindId(containerID).One(&container)
	if err != nil {
		log.Errorf("Repository Container query failed for containerID %v %v", containerID, err)
		return nil, err
	}

	h := s.DB(repo.Config.Database).C("hosts")
	host := models.DockerHost{}
	err = h.FindId(container.HostId).One(&host)
	if err != nil {
		log.Errorf("Repository Container hosts query failed for containerID %v %v", containerID, err)
		return nil, err
	}

	containers := []models.DockerContainer{}
	containers = append(containers, container)

	payload := &models.DockerPayload{
		Host:       host,
		Containers: containers,
	}

	return payload, nil
}

func (repo *MongoRepository) AllSyrosServices() ([]models.SyrosService, error) {
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

func (repo *MongoRepository) AllHealthChecks() ([]models.ConsulHealthCheck, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("checks")
	checks := []models.ConsulHealthCheck{}
	err := c.Find(nil).Sort("-collected").All(&checks)
	if err != nil {
		log.Errorf("Repository AllHealthChecks query failed %v", err)
		return nil, err
	}

	return checks, nil
}
