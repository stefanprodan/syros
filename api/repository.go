package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (repo *Repository) AllEnvironments() ([]string, error) {
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

func (repo *Repository) EnvironmentHostContainerSum() ([]models.EnvironmentStats, error) {
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

func (repo *Repository) AllHosts() ([]models.DockerHost, error) {
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

func (repo *Repository) HostContainers(hostID string) (*models.DockerPayload, error) {
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

func (repo *Repository) EnvironmentContainers(env string) (*models.DockerPayload, error) {
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
	err = c.Find(bson.M{"environment": env}).Sort("created").All(&containers)
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

func (repo *Repository) AllContainers() ([]models.DockerContainer, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("containers")
	containers := []models.DockerContainer{}
	err := c.Find(nil).Sort("-collected").All(&containers)
	if err != nil {
		log.Errorf("Repository AllContainers query failed %v", err)
		return nil, err
	}

	return containers, nil
}

func (repo *Repository) Container(containerID string) (*models.DockerPayload, error) {
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

func (repo *Repository) AllHealthChecks() ([]models.ConsulHealthCheck, error) {
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

func (repo *Repository) HealthCheckLog(checkId string) ([]models.ConsulHealthCheckLog, []models.HealthCheckStats, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("checks_log")
	logs := []models.ConsulHealthCheckLog{}
	err := c.Find(bson.M{"check_id": checkId}).Sort("-begin").Limit(500).All(&logs)
	if err != nil {
		log.Errorf("Repository HealthCheckLog checks_log query failed %v", err)
		return nil, nil, err
	}

	k := s.DB(repo.Config.Database).C("checks")
	current := models.ConsulHealthCheck{}
	err = k.FindId(checkId).One(&current)
	if err != nil {
		log.Errorf("Repository HealthCheckLog checks query failed %v", err)
		return nil, nil, err
	}

	// add current status to logs
	cur := models.NewConsulHealthCheckLog(current, current.Since, time.Now().UTC())
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
		log.Errorf("Repository HealthCheckLog pipeline failed %v", err)
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

func (repo *Repository) AllVSphere() (*models.VSpherePayload, error) {
	s := repo.Session.Copy()
	defer s.Close()

	c := s.DB(repo.Config.Database).C("vsphere_hosts")
	hosts := []models.VSphereHost{}
	err := c.Find(nil).Sort("-collected").All(&hosts)
	if err != nil {
		log.Errorf("Repository AllVSphere vsphere_hosts cursor failed %v", err)
		return nil, err
	}

	v := s.DB(repo.Config.Database).C("vsphere_vms")
	vms := []models.VSphereVM{}
	err = v.Find(nil).Sort("name").All(&vms)
	if err != nil {
		log.Errorf("Repository AllVSphere vsphere_vms cursor failed %v", err)
		return nil, err
	}

	d := s.DB(repo.Config.Database).C("vsphere_dstores")
	ds := []models.VSphereDatastore{}
	err = d.Find(nil).Sort("-collected").All(&ds)
	if err != nil {
		log.Errorf("Repository AllVSphere vsphere_dstores cursor failed %v", err)
		return nil, err
	}

	payload := &models.VSpherePayload{
		Hosts:      hosts,
		VMs:        vms,
		DataStores: ds,
	}

	return payload, nil
}
