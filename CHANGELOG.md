# Syros Changelog

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