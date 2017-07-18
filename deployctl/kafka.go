package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"bufio"
	"fmt"

	"github.com/pkg/errors"
)

type KafkaDeploy struct {
	Env     string
	HostTo  string
	Service string
	Dir     string
	Ssh     *SshClient
}

func (m KafkaDeploy) Migrate() error {
	log.Printf("Starting %s migration to %s", m.Service, m.HostTo)

	migration := make([]string, 0)

	log.Printf("Loading topics from %s", m.Dir)
	err := filepath.Walk(m.Dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".topics") {
			list, err := readTopicsCommand(path)
			if err != nil {
				return err
			}
			migration = append(migration, list...)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "loading topics from %s failed", m.Dir)
	}
	sort.Strings(migration)
	//log.Print(migration)

	log.Printf("Topics found %d in migration", len(migration))

	log.Printf("Loading topics via SSH from %s", m.HostTo)
	output, err := m.Ssh.RunCommand("/usr/bin/kafka-topics --zookeeper 127.0.0.1:2181 --list|sort")
	if err != nil {
		return errors.Wrapf(err, "loading topics via ssh from %s failed", m.HostTo)
	}

	existingTopics := make([]string, 0)
	for _, line := range strings.Split(output, "\n") {
		if len(line) > 2 {
			existingTopics = append(existingTopics, strings.TrimSpace(line))
		}
	}
	log.Printf("Topics found %d on %s", len(existingTopics), m.HostTo)
	//log.Print(existingTopics)

	newTopics := make([]string, 0)
	for _, mi := range migration {
		found := false
		for _, mr := range existingTopics {
			// strip topic specs as retention or partitions
			topicName := strings.Split(mi, " ")[0]
			if strings.TrimSpace(mr) == strings.TrimSpace(topicName) {
				found = true
				//log.Printf("Topic found %v", mi)
				break
			}
		}
		if !found {
			newTopics = append(newTopics, mi)
		}
	}

	if len(newTopics) > 0 {
		log.Printf("New topics found %v", newTopics)
		for _, topicCommand := range newTopics {
			log.Printf("Inserting new topic %s", topicCommand)
			cmd := fmt.Sprintf("/usr/bin/kafka-topics --zookeeper 127.0.0.1:2181 --create --topic %s --replication-factor 3 --partitions 1", topicCommand)
			if strings.Contains(topicCommand, "partitions") {
				cmd = fmt.Sprintf("/usr/bin/kafka-topics --zookeeper 127.0.0.1:2181 --create --topic %s --replication-factor 3", topicCommand)
			}
			log.Print(cmd)
			_, err := m.Ssh.RunCommand(cmd)
			if err != nil {
				return errors.Wrapf(err, "kafka-topics create %s failed", topicCommand)
			}
		}
	} else {
		log.Print("No new topics found, Kafka is up to date")
	}

	return nil
}

func readTopicsCommand(path string) ([]string, error) {
	topics := make([]string, 0)

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			val := scanner.Text()
			if len(val) > 2 {
				topics = append(topics, val)
			}
		}

		if err = scanner.Err(); err != nil {
			return topics, err
		}

	} else {
		return topics, err
	}

	return topics, nil
}
