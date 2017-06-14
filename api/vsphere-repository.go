package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
)

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
