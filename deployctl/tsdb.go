package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type TsdbDeploy struct {
	Env     string
	HostTo  string
	Service string
	Dir     string
	Ssh     *SshClient
}

func (m TsdbDeploy) Migrate() error {
	log.Printf("Starting %s migration to %s", m.Service, m.HostTo)

	migration := make([]string, 0)

	log.Printf("Loading metrics from %s", m.Dir)
	err := filepath.Walk(m.Dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".metrics") {
			list, err := readLines(path)
			if err != nil {
				return err
			}
			migration = append(migration, list...)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "loading metrics from %s failed", m.Dir)
	}
	migration = stringsDistinct(migration)

	log.Printf("Metrics found %d in migration", len(migration))

	log.Printf("Loading metrics via SSH from %s", m.HostTo)
	output, err := m.Ssh.RunCommand("/usr/share/opentsdb/bin/tsdb uid grep metrics .|sort")
	if err != nil {
		return errors.Wrapf(err, "loading metrics via ssh from %s failed", m.HostTo)
	}
	sort.Strings(migration)
	//log.Print(migration)

	existingMetrics := make([]string, 0)
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "metrics ") {
			metric := strings.TrimSpace(strings.Replace(strings.Split(line, ":")[0], "metrics ", "", -1))
			if len(metric) > 2 {
				existingMetrics = append(existingMetrics, metric)
			}
		}
	}
	log.Printf("Metrics found %d on %s", len(existingMetrics), m.HostTo)
	//log.Print(existingMetrics)

	newMetrics := make([]string, 0)
	for _, mi := range migration {
		found := false
		for _, mr := range existingMetrics {
			if strings.TrimSpace(mr) == strings.TrimSpace(mi) {
				found = true
				//log.Printf("Metrics found %v", mi)
				break
			}
		}
		if !found {
			newMetrics = append(newMetrics, mi)
		}
	}

	if len(newMetrics) > 0 {
		log.Printf("New metrics found %v", newMetrics)
		for _, metric := range newMetrics {
			log.Printf("Inserting new metric %s", metric)
			_, err := m.Ssh.RunCommand("/usr/share/opentsdb/bin/tsdb mkmetric " + metric)
			if err != nil {
				return errors.Wrapf(err, "mkmetric %s failed", metric)
			}
		}
	} else {
		log.Print("No new metrics found, OpenTSDB is up to date")
	}

	return nil
}

func readLines(path string) ([]string, error) {
	metrics := make([]string, 0)

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			val := scanner.Text()
			if len(val) > 2 {
				metrics = append(metrics, val)
			}
		}

		if err = scanner.Err(); err != nil {
			return metrics, err
		}

	} else {
		return metrics, err
	}

	return metrics, nil
}

func stringsDistinct(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make([]string, 0, len(s))
slice:
	for i, n := range s {
		if i == 0 {
			s = s[:0]
		}
		for _, t := range seen {
			if n == t {
				continue slice
			}
		}
		seen = append(seen, n)
		s = append(s, n)
	}
	return s
}
