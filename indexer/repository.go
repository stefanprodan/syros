package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/stefanprodan/syros/models"
	"time"
)

type Repository struct {
	Config  *Config
	Session *r.Session
}

func NewRepository(config *Config) (*Repository, error) {

	session, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDB,
		Database: config.Database,
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		Config:  config,
		Session: session,
	}
	return repo, nil
}

func (repo *Repository) Initialize() {
	var cursor *r.Cursor
	var err error
	var cnt int

	// init db
	cursor, err = r.DBList().Contains(repo.Config.Database).Run(repo.Session)
	if err != nil {
		log.Fatalf("RethinkDB database init query failed %v", err)
	}

	cursor.One(&cnt)
	cursor.Close()

	if cnt < 1 {
		log.Infof("RethinkDB no database found, creating %v", repo.Config.Database)
		_, err := r.DBCreate(repo.Config.Database).RunWrite(repo.Session)
		if err != nil {
			log.Fatalf("RethinkDB database creation failed %v", err)
		}
	}

	repo.CreateTable("hosts")
	repo.CreateIndex("hosts", "environment")
	repo.CreateIndex("hosts", "collected")
	repo.CreateTable("containers")
	repo.CreateIndex("containers", "host_id")
	repo.CreateIndex("containers", "environment")
	repo.CreateIndex("containers", "collected")
	repo.CreateTable("checks")
	repo.CreateIndex("checks", "host_id")
	repo.CreateIndex("checks", "environment")
	repo.CreateIndex("checks", "collected")
	repo.CreateTable("syros_services")
	repo.CreateIndex("syros_services", "environment")
	repo.CreateIndex("syros_services", "collected")
}

func (repo *Repository) CreateTable(table string) {
	rdb := r.DB(repo.Config.Database)
	cursor, err := rdb.TableList().Contains(table).Run(repo.Session)
	if err != nil {
		log.Fatalf("RethinkDB table init query failed %v", err)
	}
	var cnt int
	cursor.One(&cnt)
	cursor.Close()

	if cnt < 1 {
		log.Infof("RethinkDB no table found, creating %v", table)
		_, err := rdb.TableCreate(table).RunWrite(repo.Session)
		if err != nil {
			log.Fatalf("RethinkDB %v table creation failed %v", table, err)
		}
	}
}

func (repo *Repository) CreateIndex(table string, field string) {
	t := r.DB(repo.Config.Database).Table(table)
	cursor, err := t.IndexList().Contains(field).Run(repo.Session)
	if err != nil {
		log.Fatalf("RethinkDB index init query failed %v", err)
	}

	var cnt int
	cursor.One(&cnt)
	cursor.Close()

	if cnt < 1 {
		log.Infof("RethinkDB no index found on table %v, creating %v", table, field)
		err := t.IndexCreate(field).Exec(repo.Session)
		if err != nil {
			log.Fatalf("RethinkDB table %v index %v creation failed %v", table, field, err)
		}
		t.IndexWait().RunWrite(repo.Session)
		if err != nil {
			log.Fatalf("RethinkDB table %v index %v wait failed %v", table, field, err)
		}
	}
}

func (repo *Repository) HostUpsert(host models.DockerHost) {
	res, err := r.Table("hosts").Get(host.Id).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository host upsert query after ID failed %v", err)
		return
	}

	if res.IsNil() {
		_, err := r.Table("hosts").Insert(host).RunWrite(repo.Session)
		if err != nil {
			log.Errorf("Repository host insert failed %v", err)
		}
	} else {
		_, err := r.Table("hosts").Get(host.Id).Update(host).Run(repo.Session)
		if err != nil {
			log.Errorf("Repository host update failed %v", err)
		}
	}
}

func (repo *Repository) ContainerUpsert(container models.DockerContainer) {
	cursor, err := r.Table("containers").Get(container.Id).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository containers upsert query after ID failed %v", err)
		return
	}

	if cursor.IsNil() {
		_, err := r.Table("containers").Insert(container).RunWrite(repo.Session)
		if err != nil {
			log.Errorf("Repository containers insert failed %v", err)
		}
	} else {
		_, err := r.Table("containers").Get(container.Id).Update(container).Run(repo.Session)
		if err != nil {
			log.Errorf("Repository containers update failed %v", err)
		}
	}
	cursor.Close()
}

func (repo *Repository) CheckUpsert(check models.ConsulHealthCheck) {
	cursor, err := r.Table("checks").Get(check.Id).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository checks upsert query after ID failed %v", err)
		return
	}

	if cursor.IsNil() {
		check.Since = check.Collected
		_, err := r.Table("checks").Insert(check).RunWrite(repo.Session)
		if err != nil {
			log.Errorf("Repository checks insert failed %v", err)
		}
	} else {
		c := models.ConsulHealthCheck{}
		err = cursor.One(&c)
		if c.Status != check.Status {
			check.Since = check.Collected
		} else {
			check.Since = c.Since
		}

		_, err := r.Table("checks").Get(check.Id).Update(check).Run(repo.Session)
		if err != nil {
			log.Errorf("Repository checks update failed %v", err)
		}
	}
	cursor.Close()
}

func (repo *Repository) SyrosServiceUpsert(service models.SyrosService) {
	cursor, err := r.Table("syros_services").Get(service.Id).Run(repo.Session)
	if err != nil {
		log.Errorf("Repository syros_services upsert query after ID failed %v", err)
		return
	}

	if cursor.IsNil() {
		_, err := r.Table("syros_services").Insert(service).RunWrite(repo.Session)
		if err != nil {
			log.Errorf("Repository syros_services insert failed %v", err)
		}
	} else {
		_, err := r.Table("syros_services").Get(service.Id).Update(service).Run(repo.Session)
		if err != nil {
			log.Errorf("Repository syros_services update failed %v", err)
		}
	}
	cursor.Close()
}

// Removes stale records that have not been updated for a while
func (repo *Repository) RunGarbageCollector(tables []string) {
	if repo.Config.DatabaseStale > 0 {
		go func(stale int, since int) {
			log.Infof("Stating repository GC interval %v minutes", stale)
			for _, table := range tables {
				res, err := r.Table(table).
					Between(time.Now().Add(-time.Duration(since)*time.Hour).UTC(),
						time.Now().Add(-time.Duration(stale)*time.Minute).UTC(),
						r.BetweenOpts{Index: "collected"}).
					Delete().RunWrite(repo.Session)
				if err != nil {
					log.Errorf("Repository GC for table %v query failed %v", table, err)
				} else {
					if res.Deleted > 0 {
						log.Infof("Repository GC removed %v from %v", res.Deleted, table)
					}
				}
			}

			time.Sleep(60 * time.Second)
		}(repo.Config.DatabaseStale, repo.Config.DatabaseStaleSince)
	}
}

func logRepositoryResponse(action string, response interface{}) {
	jBytes, _ := json.Marshal(response)
	log.Debugf("Repository %v result %s", action, string(jBytes))
}
