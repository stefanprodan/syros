package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type VSphereCollector struct {
	ApiAddress  string
	Include     []string
	Exclude     []string
	Environment string
	Topic       string
}

func NewVSphereCollector(address string, include []string, exclude []string, env string) (*VSphereCollector, error) {

	c := &VSphereCollector{
		ApiAddress:  address,
		Include:     include,
		Exclude:     exclude,
		Environment: env,
		Topic:       "vsphere",
	}

	return c, nil
}

func (col *VSphereCollector) Collect() (*models.VSpherePayload, error) {
	start := time.Now().UTC()
	payload := &models.VSpherePayload{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := url.Parse(col.ApiAddress)
	if err != nil {
		return nil, err
	}

	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, err
	}

	pc := property.DefaultCollector(c.Client)
	f := find.NewFinder(c.Client, true)

	dc, err := f.DefaultDatacenter(ctx)
	if err != nil {
		return nil, err
	}
	f.SetDatacenter(dc)

	datastores, err := getDatastores(ctx, f, pc)
	if err != nil {
		return nil, err
	}

	clusters, err := getClusters(ctx, f, pc)
	if err != nil {
		return nil, err
	}

	hosts, err := getHosts(ctx, f, pc, clusters)
	if err != nil {
		return nil, err
	}

	vms, err := getVMs(ctx, f, pc, hosts, datastores, col.Include, col.Exclude)
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		for i, host := range hosts {
			if vm.HostId == host.Id {
				hosts[i].VMs++
			}
		}

		for i, ds := range datastores {
			if vm.DatastoreId == ds.Id {
				datastores[i].VMs++
			}
		}
	}

	payload.DataStores = datastores
	payload.Hosts = hosts
	payload.VMs = vms

	log.Debugf("%v collect duration: %v vms %v", col.ApiAddress, time.Now().UTC().Sub(start), len(payload.VMs))
	return payload, nil
}

func getDatastores(ctx context.Context, f *find.Finder, pc *property.Collector) ([]models.VSphereDatastore, error) {
	result := make([]models.VSphereDatastore, 0)
	dss, err := f.DatastoreList(ctx, "*")
	if err != nil {
		return result, err
	}

	var refs []types.ManagedObjectReference
	for _, ds := range dss {
		refs = append(refs, ds.Reference())
	}

	var dst []mo.Datastore
	err = pc.Retrieve(ctx, refs, []string{"summary", "parent"}, &dst)
	if err != nil {
		return result, err
	}

	for _, ds := range dst {
		res := models.VSphereDatastore{
			Name:      ds.Summary.Name,
			Collected: time.Now().UTC(),
			Capacity:  ds.Summary.Capacity,
			Free:      ds.Summary.FreeSpace,
			Type:      ds.Summary.Type,
			Id:        ds.Summary.Datastore.Value,
		}
		result = append(result, res)
	}

	return result, nil
}

func getClusters(ctx context.Context, f *find.Finder, pc *property.Collector) (map[string][]string, error) {
	result := make(map[string][]string, 0)

	clusters, err := f.ClusterComputeResourceList(ctx, "*")
	if err != nil {
		return result, err
	}

	var cRefs []types.ManagedObjectReference
	for _, h := range clusters {
		cRefs = append(cRefs, h.Reference())
	}

	var clusts []mo.ClusterComputeResource
	err = pc.Retrieve(ctx, cRefs, []string{"name", "host"}, &clusts)
	if err != nil {
		return result, err
	}

	for _, cl := range clusts {
		if len(cl.Host) > 0 {
			hosts := make([]string, 0)
			for _, host := range cl.Host {
				hosts = append(hosts, host.Value)
			}
			result[cl.Name] = hosts
		}
	}

	return result, nil
}

func getHostCluster(hostId string, clusters map[string][]string) string {
	for cluster, hosts := range clusters {
		for _, h := range hosts {
			if h == hostId {
				return cluster
			}
		}
	}

	return ""
}

func getHosts(ctx context.Context, f *find.Finder,
	pc *property.Collector,
	clusters map[string][]string) ([]models.VSphereHost, error) {
	result := make([]models.VSphereHost, 0)

	hosts, err := f.HostSystemList(ctx, "*")
	if err != nil {
		return result, err
	}

	var hRefs []types.ManagedObjectReference
	for _, h := range hosts {
		hRefs = append(hRefs, h.Reference())
	}

	var hostList []mo.HostSystem
	err = pc.Retrieve(ctx, hRefs, []string{"name", "summary", "hardware", "runtime"}, &hostList)
	if err != nil {
		return result, err
	}

	for _, host := range hostList {
		res := models.VSphereHost{
			Name:       host.Name,
			Id:         host.Summary.Host.Value,
			Collected:  time.Now().UTC(),
			PowerState: fmt.Sprintf("%v", host.Runtime.PowerState),
			BootTime:   host.Runtime.BootTime,
			Cluster:    getHostCluster(host.Summary.Host.Value, clusters),
			Memory:     host.Hardware.MemorySize,
			NCPU:       int(host.Hardware.CpuInfo.NumCpuPackages * host.Hardware.CpuInfo.NumCpuCores),
		}

		result = append(result, res)
	}

	return result, nil
}

func getVMs(ctx context.Context, f *find.Finder,
	pc *property.Collector,
	hosts []models.VSphereHost,
	datastores []models.VSphereDatastore,
	include []string, exclude []string) ([]models.VSphereVM, error) {
	result := make([]models.VSphereVM, 0)

	vms, err := f.VirtualMachineList(ctx, "*")
	if err != nil {
		return result, err
	}

	var vmRefs []types.ManagedObjectReference
	for _, vm := range vms {
		vmRefs = append(vmRefs, vm.Reference())
	}

	var vmt []mo.VirtualMachine
	err = pc.Retrieve(ctx, vmRefs, []string{"name", "summary", "guest", "datastore", "runtime", "storage"}, &vmt)
	if err != nil {
		return result, err
	}

	for _, vm := range vmt {
		in := applyFilter(vm.Name, include, exclude)
		if !in {
			log.Debugf("%v VM excluded by filter", vm.Name)
			continue
		}

		hasHost := false
		var host models.VSphereHost
		for _, h := range hosts {
			if h.Id == vm.Summary.Runtime.Host.Value {
				hasHost = true
				host = h
				break
			}
		}
		if !hasHost {
			log.Errorf("%v VM host not found", vm.Name)
			continue
		}

		if len(vm.Datastore) < 1 {
			log.Errorf("%v VM has no datastore", vm.Name)
			continue
		}

		hasDatastore := false
		var datastore models.VSphereDatastore
		for _, h := range datastores {
			if h.Id == vm.Datastore[0].Value {
				hasDatastore = true
				datastore = h
				break
			}
		}
		if !hasDatastore {
			log.Errorf("%v VM datastore not found", vm.Name)
			continue
		}

		res := models.VSphereVM{
			NCPU:          int(vm.Summary.Config.NumCpu),
			Memory:        int64(vm.Summary.Config.MemorySizeMB),
			Name:          vm.Name,
			Cluster:       host.Cluster,
			BootTime:      vm.Runtime.BootTime,
			PowerState:    fmt.Sprintf("%v", vm.Runtime.PowerState),
			Collected:     time.Now().UTC(),
			Id:            vm.Summary.Vm.Value,
			HostName:      host.Name,
			HostId:        host.Id,
			DatastoreId:   datastore.Id,
			DatastoreName: datastore.Name,
			Environment:   strings.Split(vm.Name, "-")[0],
		}

		if vm.Guest != nil {
			res.IP = vm.Guest.IpAddress
		}

		var storeSize int64
		for _, store := range vm.Storage.PerDatastoreUsage {
			storeSize += store.Committed + store.Uncommitted
		}
		res.Storage = storeSize

		result = append(result, res)
	}

	return result, nil
}

func applyFilter(vm string, include []string, exclude []string) bool {
	if len(exclude) > 0 {
		for _, filter := range exclude {
			if strings.Contains(vm, filter) {
				return false
			}
		}
	}

	in := false
	if len(include) > 0 {
		for _, filter := range include {
			if strings.Contains(vm, filter) {
				in = true
				break
			}
		}
	} else {
		in = true
	}

	return in
}
