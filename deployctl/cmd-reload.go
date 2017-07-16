package main

import (
	"log"
	"path"

	"github.com/urfave/cli"
)

func componentReload(c *cli.Context) error {

	config := c.GlobalString("config")
	ticket := c.String("ticket")
	environments := c.StringSlice("environment")
	components := c.StringSlice("component")
	tag := ""

	dir, err := createArtifactsDir("/tmp")
	if err != nil {
		panic(err.Error())
	}

	setLogFile(dir)

	log.Print(">>> Deployment started")
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

			if componentCfg.Component.Type == "docker" {
				targets := componentCfg.Component.Target

				// detect leader node and shift it as last in the targets array
				if componentCfg.Component.Mode == "ha" {
					leaderIdx := -1
					var leaderEntry ComponentTarget
					for i, target := range componentCfg.Component.Target {
						isLeader := containerClusterCheck(target.Leadership, 15)
						if isLeader {
							leaderIdx = i
							leaderEntry = target
							break
						}
					}
					if leaderIdx > -1 {
						targets = removeTargetByIndex(targets, leaderIdx)
						targets = append(targets, leaderEntry)
						log.Printf("Leader found on %s", leaderEntry.Host)
					} else {
						log.Printf("No leader found for %s on %s", component, env)
					}
				}

				// run reload on each target host
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

					cd := ContainerDeploy{
						Dir:      dir,
						Env:      env,
						HostFrom: target.Host,
						HostTo:   target.Host,
						Service:  component,
						Tag:      tag,
						Ticket:   ticket,
						Check:    target.Health,
					}

					err = cd.Promote()
					if err != nil {
						log.Fatal(err.Error())
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
								err := jira.Post(ticket, env, component, target.Host)
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
					}

					log.Printf("Deployment complete for %s on %s", component, target.Host)
					log.Print("-----------------")
				}
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

	log.Print(">>> Deployment complete")

	return nil
}
