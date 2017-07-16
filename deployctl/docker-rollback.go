package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func (cd ContainerDeploy) Rollback() error {
	log.Printf("Starting %s rollbackup to %s", cd.Service, cd.HostTo)

	dockerTo := fmt.Sprintf("tcp://%s:2375", cd.HostTo)

	// archive previous if exists
	if k, _ := containerExists(dockerTo, cd.Service+"-previous"); k {
		// delete current if exists
		if k, _ := containerExists(dockerTo, cd.Service); k {
			err := containerRemove(dockerTo, cd.Service)
			if err != nil {
				return errors.Wrap(err, "containerRollback failed")
			}
			log.Printf("Container %s deleted", cd.Service)
		}
		// rename previous to current
		err := containerRename(dockerTo, cd.Service+"-previous", cd.Service)
		if err != nil {
			return errors.Wrap(err, "containerRollback failed")
		}
		log.Printf("Container %s-previous renamed to %s", cd.Service, cd.Service)

		// start current
		err = containerStart(dockerTo, cd.Service)
		if err != nil {
			return errors.Wrap(err, "containerRollback failed")
		}
		log.Printf("Container %s started", cd.Service)

	} else {
		return errors.Errorf("Container %s-previous not found on %s", cd.Service, cd.HostTo)
	}

	ok, err := containerIsRunning(dockerTo, cd.Service)
	if err != nil {
		return errors.Wrap(err, "containerRollback failed")
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
