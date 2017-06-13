<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li>health</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2 class="critical">{{ stats.unhealthy }}</h2><small class="text-uppercase">Offline</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.healthy }}</h2><small class="text-uppercase">Online</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.services }}</h2><small class="text-uppercase">Nodes</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.envs }}</h2><small class="text-uppercase">Environments</small></div>
    </div>
  </div>
  <div class="charts">
    <div class="row" v-if="loaded">
      <div class="col-md-16">
        <div class="line-chart">
          <env-chart ref="envChart" :chartData="envData" :height="envHeight"></env-chart>
          <small class="text-uppercase">Clusters distribution</small>
        </div>      
      </div>
    </div>
  </div>
  <v-client-table ref="healthchecksTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>  
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/clusters/row.template.jsx'
  import rowChild from 'components/clusters/row-child.template.jsx'
  import EnvChart from 'components/clusters/env.chart.vue'

  export default {
    name: 'clusters',
    data () {
      return {
        timer: null,
        loaded: false,
        envHeight: 180,
        stats: {healthy: '0', unhealthy: '0', services: '0', envs: '0'},
        columns: ['service_name', 'status', 'since', 'host_name', 'collected'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['service_name', 'status', 'host_name', 'since', 'collected'],
          dateColumns: ['since', 'collected'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'host_name', ascending: true},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          childRow: rowChild,
          templates: rowTemplate
        }
      }
    },
    components: {
      EnvChart
    },
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get('/cluster/healthchecks')
          .then((response) => {
            if (response != null) {
              this.tableData = response.data
              var statsHealthy = 0
              var statsUnhealthy = 0
              var statsServices = 0
              var statsEnvs = []
              for (var i = 0, len = response.data.length; i < len; i++) {
                statsServices += 1
                if (response.data[i].status !== 'offline') {
                  statsHealthy += 1
                } else {
                  statsUnhealthy += 1
                }
                statsEnvs.push(response.data[i].environment)
              }
              let envs = [...new Set(statsEnvs)]
              this.stats = {
                healthy: statsHealthy.toString(),
                unhealthy: statsUnhealthy.toString(),
                services: statsServices.toString(),
                envs: envs.length
              }

              this.envData = this.fillChart(response.data, envs.sort())
              this.loaded = true

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
      fillChart (data, envs) {
        var passingData = []
        var criticalData = []
        for (var e = 0, elen = envs.length; e < elen; e++) {
          var critical = 0
          var passing = 0
          for (var i = 0, len = data.length; i < len; i++) {
            if (data[i].environment === envs[e]) {
              if (data[i].status === 'leader') {
                passing += 1
              } else {
                critical += 1
              }
            }
          }
          passingData.push(passing)
          criticalData.push(critical)
        }

        return {
          labels: envs,
          datasets: [ {
            label: 'leaders',
            backgroundColor: '#2ECC71',
            borderWidth: 0,
            data: passingData,
            stack: '1'
          }, {
            label: 'followers',
            backgroundColor: '#F7DC6F',
            borderWidth: 0,
            data: criticalData,
            stack: '1'
          }]
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
