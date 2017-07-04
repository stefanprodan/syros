package main

import (
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/pkg/errors"
)

type ContainerDeploy struct {
	Ticket   string
	Env      string
	HostTo   string
	HostFrom string
	Tag      string
	Service  string
	Dir      string
	Check    string
}

func (cd ContainerDeploy) Promote() error {
	log.Printf("Starting %s promotion to %s", cd.Service, cd.HostTo)

	dockerTo := fmt.Sprintf("tcp://%s:2375", cd.HostTo)

	// build compose for target host
	target := strings.Split(cd.HostTo, ".")[0]
	cdir, err := buildCompose(cd.Dir, cd.Env, target)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}
	log.Printf("Docker Compose built for %s", target)

	// get vNext tag
	if len(cd.Tag) < 1 {
		dockerFrom := fmt.Sprintf("tcp://%s:2375", cd.HostFrom)
		cd.Tag, err = imageGetTag(dockerFrom, cd.Service)
		if err != nil {
			return errors.Wrap(err, "containerPromote failed")
		}
		log.Printf("Tag acquired from %s", cd.HostFrom)
	}
	log.Printf("Tag to be deployed %s", cd.Tag)

	// pull vNext on target host
	err = imagePull(dockerTo, cdir, cd.Service, cd.Tag)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}
	log.Printf("Image %s:%s pulled on %s", cd.Service, cd.Tag, cd.HostTo)

	// archive previous if exists
	if k, _ := containerExists(dockerTo, cd.Service+"-previous"); k {
		// delete old purge if exists
		if k, _ := containerExists(dockerTo, cd.Service+"-purge"); k {
			err = containerRemove(dockerTo, cd.Service+"-purge")
			if err != nil {
				return errors.Wrap(err, "containerPromote failed")
			}
			log.Printf("Container %s-purge deleted", cd.Service)
		}
		err = containerRename(dockerTo, cd.Service+"-previous", cd.Service+"-purge")
		if err != nil {
			return errors.Wrap(err, "containerPromote failed")
		}
		log.Printf("Container %s-previous renamed to purge", cd.Service)
	} else {
		log.Printf("Container %s-previous not found", cd.Service)
	}

	// set current as previous
	if k, _ := containerExists(dockerTo, cd.Service); k {
		err = containerRename(dockerTo, cd.Service, cd.Service+"-previous")
		if err != nil {
			return errors.Wrap(err, "containerPromote failed")
		}
		log.Printf("Current container renamed to %s-previous", cd.Service)
	}

	// create vNext
	err = containerCreate(dockerTo, cdir, cd.Service, cd.Tag, cd.Ticket)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}
	log.Printf("Container %s:%s created", cd.Service, cd.Tag)

	// stop previous
	if k, _ := containerExists(dockerTo, cd.Service+"-previous"); k {
		err = containerStop(dockerTo, cd.Service+"-previous")
		if err != nil {
			return errors.Wrap(err, "containerPromote failed")
		}
		log.Print("Previous container stopped")
	}

	// start vNext
	err = containerStart(dockerTo, cd.Service)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}
	log.Printf("Container %s:%s started", cd.Service, cd.Tag)

	// delete purge
	if k, _ := containerExists(dockerTo, cd.Service+"-purge"); k {
		err = containerRemove(dockerTo, cd.Service+"-purge")
		if err != nil {
			return errors.Wrap(err, "containerPromote failed")
		}
		log.Printf("Container %s-purge deleted", cd.Service)
	}

	err = imagePurge(dockerTo)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}
	log.Print("Docker images cleanup done")

	ok, err := containerIsRunning(dockerTo, cd.Service)
	if err != nil {
		return errors.Wrap(err, "containerPromote failed")
	}

	if ok {
		log.Printf("Container %s is running", cd.Service)
	} else {
		return errors.Errorf("Container %s is not running", cd.Service)
	}

	if len(cd.Check) > 0 {
		log.Printf("Begining health check for %s", cd.Check)
		ok = false
		for i := 0; i < 10; i++ {
			log.Printf("Checking health try %d", i+1)
			time.Sleep(10 * time.Second)
			err = containerHealthCheck(cd.Check, 15)
			if err != nil {
				log.Print(err.Error())
			} else {
				ok = true
				break
			}
		}

		if ok {
			log.Printf("Container %s is healthy", cd.Service)
		} else {
			return errors.Errorf("Container %s is not healthy", cd.Service)
		}
	}

	return nil
}
