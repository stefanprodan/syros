<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li>vsphere</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.vms }}</h2><small class="text-uppercase">VMS</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.vcpus }}</h2><small class="text-uppercase">vCPUs</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.vram }}</h2><small class="text-uppercase">Memory</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.vdisk }}</h2><small class="text-uppercase">Storage</small></div>
    </div>
  </div>
  <div class="charts">
    <div class="row" v-if="loaded">
      <div class="col-md-16">
        <div class="line-chart">
          <vm-chart ref="vmChart" :chartData="vmChartData" :height="chartHeight"></vm-chart>
          <small class="text-uppercase">VMs Distribution</small>
        </div>      
      </div>
    </div>
  </div>
  <v-client-table ref="vmsTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>  
  <div class="stats">
    <div class="row">
      <div class="col-md-4 text-center">
        <h2>{{ stats.dstores }}</h2><small class="text-uppercase">Datastores</small></div>
      <div class="col-md-4 text-center">
        <h2>{{ stats.dcap }}</h2><small class="text-uppercase">capacity</small></div>
      <div class="col-md-4 text-center">
        <h2>{{ stats.dfree }}</h2><small class="text-uppercase">free</small></div>
    </div>
  </div>
  <v-client-table ref="storesTabel" :data="storesData" :columns="storesColumns" :options="storesOptions"></v-client-table>
  <div class="stats">
    <div class="row">
      <div class="col-md-4 text-center">
        <h2>{{ stats.hosts }}</h2><small class="text-uppercase">Hosts</small></div>
      <div class="col-md-4 text-center">
        <h2>{{ stats.hthreads }}</h2><small class="text-uppercase">CPU Threads</small></div>
      <div class="col-md-4 text-center">
        <h2>{{ stats.hram }}</h2><small class="text-uppercase">Memory</small></div>
    </div>
  </div>
  <v-client-table ref="hostsTabel" :data="hostsData" :columns="hostsColumns" :options="hostsOptions"></v-client-table> 
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/vsphere/row.template.jsx'
  import rowChild from 'components/vsphere/row-child.template.jsx'
  import rowTemplateStore from 'components/vsphere/row.template.ds.jsx'
  import rowTemplateHost from 'components/vsphere/row.template.hosts.jsx'
  import VmChart from 'components/vsphere/vm.chart.vue'

  export default {
    name: 'vsphere',
    data () {
      return {
        timer: null,
        loaded: false,
        chartHeight: 180,
        vmChartData: null,
        stats: {vms: '0', vcpus: '0', vram: '0 MB', vdisk: '0 MB', hosts: '0', hthreads: '0', hram: '0 MB', dstores: '0', dcap: '0 MB', dfree: '0 MB'},
        columns: ['name', 'power_state', 'ip', 'host_name', 'ncpu', 'memory', 'storage', 'boot_time'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['name', 'power_state', 'ip', 'host_name', 'ncpu', 'memory', 'storage', 'boot_time'],
          dateColumns: ['boot_time'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'name', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          childRow: rowChild,
          templates: rowTemplate
        },
        storesColumns: ['n', 'name', 'power_state', 'vms', 'type', 'capacity', 'free', 'collected'],
        storesData: [],
        storesOptions: {
          skin: 'table-hover',
          sortable: ['name', 'type', 'vms', 'capacity', 'free', 'collected'],
          dateColumns: ['collected'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'name', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          templates: rowTemplateStore
        },
        hostsColumns: ['n', 'name', 'power_state', 'vms', 'cluster', 'ncpu', 'memory', 'boot_time'],
        hostsData: [],
        hostsOptions: {
          skin: 'table-hover',
          sortable: ['name', 'power_state', 'vms', 'cluster', 'ncpu', 'memory', 'boot_time'],
          dateColumns: ['boot_time'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'name', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          templates: rowTemplateHost
        }
      }
    },
    components: {
      VmChart
    },
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get('/vsphere')
          .then((response) => {
            if (response != null) {
              this.tableData = response.data.vms
              this.storesData = response.data.data_stores
              this.hostsData = response.data.hosts
              this.vmChartData = this.fillChart(response.data.chart)
              this.loaded = true
              var statsVcpus = 0
              var statsVdisk = 0
              var statsVram = 0
              var statsHthreads = 0
              var statsHram = 0
              var statsDcap = 0
              var statsDfree = 0
              for (var i = 0, lv = response.data.vms.length; i < lv; i++) {
                statsVcpus += response.data.vms[i].ncpu
                statsVdisk += response.data.vms[i].storage
                statsVram += response.data.vms[i].memory
              }
              for (var y = 0, lh = response.data.hosts.length; y < lh; y++) {
                statsHthreads += response.data.hosts[y].ncpu
                statsHram += response.data.hosts[y].memory
              }
              for (var z = 0, ld = response.data.data_stores.length; z < ld; z++) {
                statsDcap += response.data.data_stores[z].capacity
                statsDfree += response.data.data_stores[z].free
              }
              this.stats = {
                hosts: response.data.hosts.length.toString(),
                vms: response.data.vms.length.toString(),
                vcpus: statsVcpus.toString(),
                vram: parseInt(parseFloat((statsVram / Math.pow(1024, 2))).toFixed(0)) + 'TB',
                vdisk: parseInt(parseFloat((statsVdisk / Math.pow(1024, 4))).toFixed(0)) + 'TB',
                hthreads: statsHthreads.toString(),
                hram: parseInt(parseFloat((statsHram / Math.pow(1024, 4))).toFixed(0)) + 'TB',
                dstores: response.data.data_stores.length.toString(),
                dcap: parseInt(parseFloat((statsDcap / Math.pow(1024, 4))).toFixed(0)) + 'TB',
                dfree: parseInt(parseFloat((statsDfree / Math.pow(1024, 4))).toFixed(0)) + 'TB'
              }
              this.$Progress.finish()
            } else {
              this.$Progress.fail()
            }
          })
          .catch((error) => {
            if (!error.response) {
              bus.$emit('flashMessage', {
                type: 'warning',
                message: 'Network error! Could not connect to the server'
              })
            } else {
              bus.$emit('flashMessage', {
                type: 'warning',
                message: `${error.response.statusText}! ${error.response.data}`
              })
            }
            this.$Progress.fail()
          })
      },
      refreshData () {
        this.loadData()
        console.log('Refresh data: ' + this.$options.name)
        // enqueue new call after 60 seconds
        if (this.timer) clearTimeout(this.timer)
        this.timer = setTimeout(this.refreshData, 60000)
      },
      fillChart (data) {
        var backgroundColors = []
        for (var i = 0, len = data.labels.length; i < len; i++) {
          backgroundColors.push('#309292')
        }
        return {
          labels: data.labels,
          datasets: [
            {
              label: 'vms',
              backgroundColor: backgroundColors,
              borderWidth: 0,
              data: data.values
            }
          ]
        }
      }
    },
    watch: {
      '$route' (to, from) {
        if (from.params.id !== to.params.id) {
          this.id = to.params.id
          return this.refreshData()
        }
      }
    },
    created: function () {
      console.log('Created: ' + this.$options.name)
    },
    mounted: function () {
      console.log('Mounted: ' + this.$options.name)
      this.refreshData()
    },
    destroyed: function () {
      if (this.timer) {
        clearTimeout(this.timer)
        console.log('Destroyed: ' + this.$options.name)
      }
    }
}

</script>
