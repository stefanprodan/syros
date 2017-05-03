# Syros Changelog

### v0.5.0 (release)

* Agent instrumentation with Prometheus
* Run collectors with go-cron

### v0.4.0 (release)

* API instrumentation with Prometheus
* Indexer instrumentation with Prometheus
* Collectors refactoring
* Setup http pprof

### v0.3.0 (release)

***Release tracker***

* Expose web API to record releases
* Releases history page (stats, graph and master-detail view)
* Released services page

### v0.2.0 (release)

* Record Consul health status changes
* Service health check history page 
* Service health check stats for the last 30 days

### v0.1.0 (release)

* replace RethinkDB with MongoDB
* GC refactoring
* Environments stats query optimization

### v0.0.2 (release)

Features:

***Agent***

* Collect health checks via Consul API
* Extract gliderlabs/registrator meta from container env vars to determine Consul service names

***Indexer***

* Consul service health check aggregation and db persistence 

***API***

* Consul service health check db repository and REST endpoint

***UI***

* Health dashboard (stats, graph and master-detail view)

### v0.0.1 (release)

Features:

***Agent***

* Collect Docker host information via Docker engine API
* Collect Docker containers information via Docker engine API
* Service registry via NATS 

***Indexer***

* Docker hosts and containers aggregation and db persistence 
* Agents service registry via NATS with db persistence 
* Indexer service registry db persistence
* DB migration on startup 
* DB garbage collector

***API***

* JWT auth
* Docker hosts and containers db repository 
* Docker hosts and containers REST endpoints

***UI***

* JWT auth
* Home dashboard (environments stats and listing)
* Hosts dashboard (hosts stats and master-detail view)
* Environment dashboard (environment stats and containers master-detail view)
* Container page (container stats, props, labels, env variables)
* Admin dashboard (Syros stats, agents and indexers registry master-detail view)