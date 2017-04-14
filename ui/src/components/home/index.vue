<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li>Home</li>
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
      <div class="col-md-6">
        <div class="pie-chart">
          <env-pie-chart ref="hostChart" :chartData="hostChart"></env-pie-chart>
          <small class="text-uppercase">Hosts distribution</small>
        </div>      
      </div>
      <div class="col-md-6">
        <div class="pie-chart">
          <env-pie-chart ref="containerChart" :chartData="containerChart"></env-pie-chart>
          <small class="text-uppercase">Containers distribution</small>
        </div>        
      </div>
    </div>
  </div>
  <v-client-table ref="envsTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/home/row.template.jsx'
  import EnvPieChart from 'components/home/env.piechart.vue'

  export default {
    name: 'home',
    data () {
      return {
        timer: null,
        loaded: false,
        hostChart: null,
        containerChart: null,
        stats: {hosts: '', containers: '0', cpus: '0', ram: '0 MB'},
        columns: ['n', 'environment', 'hosts', 'containers_running', 'ncpu', 'mem_total'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['environment', 'hosts', 'containers_running', 'ncpu', 'mem_total'],
          uniqueKey: 'environment',
          orderBy: {column: 'environment', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          templates: rowTemplate
        }
      }
    },
    components: {
      EnvPieChart
    },
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get('/docker/environments/stats')
          .then((response) => {
            if (response != null) {
              this.tableData = response.data
              var statsHosts = 0
              var statsContainers = 0
              var statsCpus = 0
              var statsRam = 0
              var labels = []
              var hostdata = []
              var containerdata = []
              for (var i = 0, len = response.data.length; i < len; i++) {
                statsHosts += response.data[i].hosts
                statsContainers += response.data[i].containers_running
                statsCpus += response.data[i].ncpu
                statsRam += parseInt(parseFloat((response.data[i].mem_total / Math.pow(1024, 3))).toFixed(0))

                labels.push(response.data[i].environment)
                hostdata.push(response.data[i].hosts)
                containerdata.push(response.data[i].containers_running)
              }

              this.hostChart = this.fillChart(labels, hostdata)
              this.containerChart = this.fillChart(labels, containerdata)
              this.loaded = true

              this.stats = {
                hosts: statsHosts.toString(),
                containers: statsContainers.toString(),
                cpus: statsCpus.toString(),
                ram: statsRam.toString() + 'GB'
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
        // enqueue new call after 30 seconds
        if (this.timer) clearTimeout(this.timer)
        this.timer = setTimeout(this.refreshData, 30000)
      },
      fillChart (labels, data) {
        return {
          labels: labels,
          datasets: [
            {
              backgroundColor: [
                'rgba(65, 184, 131, .8)',
                'rgba(228, 102, 81, .8)',
                'rgba(0, 116, 255, .8)',
                'rgba(155, 89, 182, .8)',
                'rgba(88, 172, 11, .8)',
                'rgba(65, 90, 131, .8)',
                'rgba(0, 216, 255, .8)'
              ],
              borderWidth: 0,
              data: data
            }
          ]
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

