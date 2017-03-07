<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li>{{ id }}</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.hosts }}</h2><small class="text-uppercase">Hosts</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.containers }}</h2><small class="text-uppercase">Containers up</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.cpus }}</h2><small class="text-uppercase">vCPUs</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.ram }}</h2><small class="text-uppercase">Memory</small></div>
    </div>
  </div>
  <div class="charts">
    <div class="row" v-if="loaded">
      <div class="col-md-16">
        <div class="line-chart">
          <deployment-chart ref="deploymentChart" :chartData="deploymentData" :height="deploymentHeight"></deployment-chart>
          <small class="text-uppercase">Deployments</small>
        </div>      
      </div>
    </div>
  </div>
  <v-client-table ref="containersTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/environment/row.template.jsx'
  import rowChild from 'components/environment/row-child.template.jsx'
  import DeploymentChart from 'components/environment/deployment.chart.vue'

  export default {
    name: 'environment',
    data () {
      return {
        timer: null,
        loaded: false,
        deploymentHeight: 180,
        deploymentData: null,
        id: this.$route.params.id,
        stats: {hosts: '', containers: '0', cpus: '0', ram: '0 MB'},
        columns: ['name', 'state', 'status', 'host_name', 'network_mode', 'port', 'created'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['name', 'state', 'status', 'host_name', 'network_mode', 'port', 'created'],
          dateColumns: ['created'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'name', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          childRow: rowChild,
          templates: rowTemplate
        }
      }
    },
    components: {
      DeploymentChart
    },
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get(`/docker/environments/${this.id}`)
          .then((response) => {
            if (response != null) {
              this.tableData = response.data.containers
              this.deploymentData = this.fillChart(response.data.deployments)
              this.loaded = true
              console.log(response.data.deployments)
              this.stats = {
                hosts: response.data.host.containers.toString(),
                containers: response.data.host.containers_running.toString(),
                cpus: response.data.host.ncpu.toString(),
                ram: parseInt(parseFloat((response.data.host.mem_total / Math.pow(1024, 3))).toFixed(0)).toString() + 'GB'
              }
              this.$Progress.finish()
            } else {
              this.$Progress.fail()
            }
          })
          .catch((error) => {
            if (!error.response.status) {
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
        // enqueue new call after 30 seconds
        if (this.timer) clearTimeout(this.timer)
        this.timer = setTimeout(this.refreshData, 30000)
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
              label: 'deployments',
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
