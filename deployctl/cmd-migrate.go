package main

import (
	"fmt"
	"log"
	"path"

	"github.com/urfave/cli"
)

func componentMigrate(c *cli.Context) error {

	config := c.GlobalString("config")
	ticket := c.String("ticket")
	environments := c.StringSlice("environment")
	components := c.StringSlice("component")

	dir, err := createArtifactsDir("/tmp")
	if err != nil {
		panic(err.Error())
	}

	setLogFile(dir)

	log.Print(">>> Migration started")
	err = downloadArtifacts(config, dir)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Config downloaded to %s", dir)
	log.Print("-----------------")

	for _, component := range components {
		for _, env := range environments {

			componentCfg, exists, err := loadComponent(dir, env, component)
			if err != nil {
				log.Fatal(err.Error())
			}
			if !exists {
				log.Printf("No config found for %s on %s", component, env)
				log.Print("-----------------")
				continue
			}

			targets := componentCfg.Component.Target

			// run migration on each target host
			for _, target := range targets {
				// mark deployment as started in SYROS
				if len(ticket) > 0 {
					syrosApi, cfgExists, err := loadSyrosConfig(dir, "syros")
					if err != nil {
						log.Printf("Syros config load failed %s", err.Error())
					} else {
						if !cfgExists {
							log.Print("Syros config not found")
						} else {
							err := syrosApi.Start(ticket, env, component, target.Host)
							if err != nil {
								log.Print(err.Error())
							}
						}
					}
				}

				// download migrations from Jenkins
				jenkinsConfig, cfgExists, err := loadJenkinsConfig(dir, "jenkins")
				if err != nil {
					log.Fatal(err.Error())
				}
				if !cfgExists {
					log.Fatal("Jenkins config not found")
				}
				url := fmt.Sprintf("%s/job/%s/lastSuccessfulBuild/artifact/packaging/%s.tar.gz",
					jenkinsConfig.API.URL, component, component)
				log.Printf("Downloading migration from %s", url)
				err = downloadArtifacts(url, dir)
				if err != nil {
					log.Fatal(err.Error())
				}

				// load ssh identity
				sshConfig, cfgExists, err := loadSshConfig(dir, "ssh")
				if err != nil {
					log.Fatal(err.Error())
				}
				if !cfgExists {
					log.Fatal("SSH config not found")
				}

				log.Printf("Setup SSH for %s@%s", sshConfig.User, target.Host)
				ssh, err := NewSshClient(sshConfig.User, target.Host, 22, sshConfig.Key, "")
				if err != nil {
					log.Fatal(err.Error())
				}

				if componentCfg.Component.Type == "tsdbmetrics" {
					cd := TsdbDeploy{
						Dir:     path.Join(dir, component, "metrics"),
						Env:     env,
						HostTo:  target.Host,
						Service: component,
						Ssh:     ssh,
					}

					err = cd.Migrate()
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				if componentCfg.Component.Type == "kafkatopics" {
					cd := KafkaDeploy{
						Dir:     path.Join(dir, component, "topics"),
						Env:     env,
						HostTo:  target.Host,
						Service: component,
						Ssh:     ssh,
					}

					err = cd.Migrate()
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				if componentCfg.Component.Type == "postgres" {
					url := fmt.Sprintf("%s/job/%s/lastSuccessfulBuild/artifact/%s.tar.gz",
						jenkinsConfig.API.URL, "flyway", "flyway")
					log.Printf("Downloading migration tool from %s", url)
					err = downloadArtifacts(url, dir)
					if err != nil {
						log.Fatal(err.Error())
					}

					cd := PostgresDeploy{
						Dir:      dir,
						Env:      env,
						HostTo:   target.Host,
						Service:  component,
						Ssh:      ssh,
						Url:      target.URL,
						User:     target.DBUser,
						Password: target.DBPassword,
						Database: target.DBName,
						Location: path.Join(dir, component, "migration"),
					}

					err = cd.Migrate()
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				if len(ticket) > 0 {
					// add comment on JIRA ticket
					jira, cfgExists, err := loadJiraConfig(dir, "jira")
					if err != nil {
						log.Printf("Jira config load failed %s", err.Error())
					} else {
						if !cfgExists {
							log.Print("Jira config not found")
						} else {
							err := jira.Post(ticket, "Migration", env, component, target.Host)
							if err != nil {
								log.Print(err.Error())
							}
						}
					}
					// mark deployment as done in SYROS and upload log
					syrosApi, cfgExists, err := loadSyrosConfig(dir, "syros")
					if err != nil {
						log.Printf("Syros config load failed %s", err.Error())
					} else {
						if !cfgExists {
							log.Print("Syros config not found")
						} else {
							err := syrosApi.Finish(ticket, env, component, target.Host, path.Join(dir, "deployctl.log"))
							if err != nil {
								log.Print(err.Error())
							}
						}
					}
					// add comment on Slack
					slack, cfgExists, err := loadSlackConfig(dir, "slack")
					if err != nil {
						log.Printf("Slack config load failed %s", err.Error())
					} else {
						if !cfgExists {
							log.Print("Slack config not found")
						} else {
							err := slack.Post(ticket, "Migration", env, component, target.Host)
							if err != nil {
								log.Print(err.Error())
							}
						}
					}
				}

				log.Printf("Migration complete for %s on %s", component, target.Host)
				log.Print("-----------------")
			}
		}
	}

	// upload log to JIRA
	if len(ticket) > 0 {
		jira, cfgExists, err := loadJiraConfig(dir, "jira")
		if err != nil {
			log.Printf("Jira config load failed %s", err.Error())
		} else {
			if !cfgExists {
				log.Print("Jira config not found")
			} else {
				err := jira.Upload(ticket, dir, "deployctl.log")
				if err != nil {
					log.Print(err.Error())
				}
			}
		}
	}

	log.Print(">>> Migration complete")

	return nil
}
