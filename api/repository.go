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

func (repo *Repository) AllEnvironments() ([]string, error) {
	cursor, err := r.Table("hosts").Distinct(r.DistinctOpts{Index: "environment"}).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository AllEnvironments query failed %v", err)
		return nil, err
	}

	environments := []string{}
	err = cursor.All(&environments)
	if err != nil {
		log.Errorf("Repository AllEnvironments cursor failed %v", err)
		return nil, err
	}
	cursor.Close()

	return environments, nil
}

func (repo *Repository) AllHosts() ([]models.DockerHost, error) {
	cursor, err := r.Table("hosts").OrderBy(r.Asc("collected"), r.OrderByOpts{Index: "collected"}).Run(repo.Session)
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

func (repo *Repository) HostContainers(hostID string) (*models.DockerPayload, error) {
	cursor, err := r.Table("hosts").Get(hostID).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository HostContainers query failed for hostID %v %v", hostID, err)
		return nil, err
	}
	host := models.DockerHost{}
	err = cursor.One(&host)
	if err != nil {
		log.Errorf("Repository HostContainers cursor failed for hostID %v %v", hostID, err)
		return nil, err
	}
	cursor.Close()

	cursor, err = r.Table("containers").GetAllByIndex("host_id", hostID).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository HostContainers query containers GetAllByIndex for hostID %v failed %v", hostID, err)
		return nil, err
	}

	containers := []models.DockerContainer{}
	err = cursor.All(&containers)
	if err != nil {
		log.Errorf("Repository HostContainers cursor containers GetAllByIndex for hostID %v failed %v", hostID, err)
		return nil, err
	}
	cursor.Close()

	payload := &models.DockerPayload{
		Host:       host,
		Containers: containers,
	}

	return payload, nil
}

func (repo *Repository) EnvironmentContainers(env string) (*models.DockerPayload, error) {
	cursor, err := r.Table("hosts").GetAllByIndex("environment", env).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers query containers GetAllByIndex for env %v failed %v", env, err)
		return nil, err
	}

	hosts := []models.DockerHost{}
	err = cursor.All(&hosts)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers cursor containers GetAllByIndex for env %v failed %v", env, err)
		return nil, err
	}
	cursor.Close()

	envStats := models.DockerHost{}

	for _, host := range hosts {
		envStats.ContainersRunning += host.ContainersRunning
		envStats.Containers++
		envStats.NCPU += host.NCPU
		envStats.MemTotal += host.MemTotal
	}

	cursor, err = r.Table("containers").GetAllByIndex("environment", env).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers query containers GetAllByIndex for env %v failed %v", env, err)
		return nil, err
	}

	containers := []models.DockerContainer{}
	err = cursor.All(&containers)
	if err != nil {
		log.Errorf("Repository EnvironmentContainers cursor containers GetAllByIndex for env %v failed %v", env, err)
		return nil, err
	}
	cursor.Close()

	payload := &models.DockerPayload{
		Host:       envStats,
		Containers: containers,
	}

	return payload, nil
}

func (repo *Repository) AllContainers() ([]models.DockerContainer, error) {
	cursor, err := r.Table("containers").OrderBy(r.Asc("collected"), r.OrderByOpts{Index: "collected"}).Run(repo.Session)
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

func (repo *Repository) Container(containerID string) (*models.DockerPayload, error) {
	cursor, err := r.Table("containers").Get(containerID).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository Container query failed for containerID %v %v", containerID, err)
		return nil, err
	}
	container := models.DockerContainer{}
	err = cursor.One(&container)
	if err != nil {
		log.Errorf("Repository Container cursor failed for containerID %v %v", containerID, err)
		return nil, err
	}
	cursor.Close()

	cursor, err = r.Table("hosts").Get(container.HostId).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository Container hosts query failed for containerID %v %v", containerID, err)
		return nil, err
	}
	host := models.DockerHost{}
	err = cursor.One(&host)
	if err != nil {
		log.Errorf("Repository Container hosts cursor failed for containerID %v %v", containerID, err)
		return nil, err
	}
	cursor.Close()

	containers := []models.DockerContainer{}
	containers = append(containers, container)

	payload := &models.DockerPayload{
		Host:       host,
		Containers: containers,
	}

	return payload, nil
}
