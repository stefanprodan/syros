# Syros deployctl

Deploy tool for Docker containers

### Prerequisite

* Docker >= 17.05
* Docker Compose >= 1.14
* curl
* tar

### Pipelines

* Docker container promotion from one env to another
* Rolling update of HA clusters

### Integrations

* JIRA ticket update
* JIRA deploy log upload
* SYROS releases update 

### Install

Latest stable version:

```bash
SYROS_VERSION=$(curl -s -o /dev/null -I -w "%{redirect_url}\n" https://github.com/stefanprodan/syros/releases/latest | grep -oP "[0-9]+(\.[0-9]+)+$")
curl -o /usr/local/bin/syros-deployctl -L https://github.com/stefanprodan/syros/releases/download/$SYROS_VERSION/syros-deployctl
chmod +x /usr/local/bin/syros-deployctl

syros-deployctl -h

```

### Usage

```bash
$ syros-deployctl -h
NAME:
   deployctl - SYROS deploy CLI

USAGE:
   syros-deployctl [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Stefan Prodan

COMMANDS:
     promote  Promote containers from one environment to another
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Download URL for the config.tar.gz file [$DCTL_CONFIG_URL]
   --help, -h                show help
   --version, -v             print the version
```

***promote***

```bash
$ syros-deployctl promote -h
NAME:
   syros-deployctl promote - Promote containers from one environment to another

USAGE:
   syros-deployctl promote [command options] [arguments...]

OPTIONS:
   --ticket value, -t value       JIRA ticket ID, if specified the deploy log will be posted on the ticket
   --environment value, -e value  Target environment, multiple values accepted
   --component value, -c value    Docker service, multiple values accepted
   --tag value                    If a tag is specified this exact docker image tag will be deployed
```

***reload***

```bash
$ syros-deployctl reload -h
NAME:
   syros-deployctl reload - Reload containers configuration

USAGE:
   syros-deployctl reload [command options] [arguments...]

OPTIONS:
   --ticket value, -t value       JIRA ticket ID, if specified the deploy log will be posted on the ticket
   --environment value, -e value  Target environment, multiple values accepted
   --component value, -c value    Docker service, multiple values accepted
```
