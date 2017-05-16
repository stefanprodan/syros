package models

import "time"

type VSpherePayload struct {
	Hosts      []VSphereHost      `json:"hosts"`
	DataStores []VSphereDatastore `json:"data_stores"`
	VMs        []VSphereVM        `json:"vms"`
}

type VSphereDatastore struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Type        string    `bson:"type" json:"type"`
	Capacity    int64     `bson:"capacity" json:"capacity"`
	Free        int64     `bson:"free" json:"free"`
	Collected   time.Time `bson:"collected" json:"collected"`
	Environment string    `bson:"environment" json:"environment"`
}

type VSphereHost struct {
	Id          string     `bson:"_id,omitempty" json:"id"`
	Name        string     `bson:"name" json:"name"`
	Cluster     string     `bson:"cluster" json:"cluster"`
	PowerState  string     `bson:"power_state" json:"power_state"`
	BootTime    *time.Time `bson:"boot_time" json:"boot_time"`
	NCPU        int        `bson:"ncpu" json:"ncpu"`
	Memory      int64      `bson:"memory" json:"memory"`
	Collected   time.Time  `bson:"collected" json:"collected"`
	Environment string     `bson:"environment" json:"environment"`
}

type VSphereVM struct {
	Id            string     `bson:"_id,omitempty" json:"id"`
	HostId        string     `bson:"host_id" json:"host_id"`
	HostName      string     `bson:"host_name" json:"host_name"`
	Cluster       string     `bson:"cluster" json:"cluster"`
	DatastoreId   string     `bson:"datastore_id" json:"datastore_id"`
	DatastoreName string     `bson:"datastore_name" json:"datastore_name"`
	Name          string     `bson:"name" json:"name"`
	PowerState    string     `bson:"power_state" json:"power_state"`
	BootTime      *time.Time `bson:"boot_time" json:"boot_time"`
	NCPU          int        `bson:"ncpu" json:"ncpu"`
	Memory        int64      `bson:"memory" json:"memory"`
	Storage       int64      `bson:"mem_total" json:"mem_total"`
	IP            string     `json:"ip"`
	Collected     time.Time  `bson:"collected" json:"collected"`
	Environment   string     `bson:"environment" json:"environment"`
}
