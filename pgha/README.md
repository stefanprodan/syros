# Syros PGHA

PostgreSQL HA agent

### Features

* Leader election via Consul
* Failover automation via repmgr
* Replication monitoring
* Health checks
* Reverse proxy via HAProxy

### Web API

* `/status` used by HAProxy to determine the current PostgreSQL master node
* `/fallback` triggers a failover if current node is master
* `/metrics` exports metrics for Prometheus scrapping
* `/config` retries node configuration

### Configuration

```
$ ./syros-pgha -h

Usage of ./syros-pgha:
  -ConsulKV string
        Consul KV prefix (default "pgha")
  -ConsulRetry int
        Number of Consul connection reties (default 10)
  -ConsulTTL string
        Consul session TTL (default "10s")
  -ConsulURI string
        Consul address (default "localhost:8500")
  -Environment string
        Environment dev|int|stg|test|prep|prod (default "dev")
  -Hostname string
        Hostname
  -LogLevel string
        logging threshold level: debug|info|warn|error|fatal|panic (default "debug")
  -NatsURI string
        Nats URI (default "nats://localhost:4222")
  -Port int
        HTTP port to listen on (default 9898)
  -PostgresCheck int
        Postgres checks interval in seconds (default 5)
  -PostgresURI string
        Postgres URI (default "postgres://user:password@localhost/db?sslmode=disable")
  -User string
        User to run under (default "postgres")
```

