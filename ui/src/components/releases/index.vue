<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li>releases</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-6 text-center">
        <h2>{{ stats.rels }}</h2><small class="text-uppercase">Releases</small></div>
      <div class="col-md-6 text-center">
        <h2>{{ stats.deps }}</h2><small class="text-uppercase">Deployments</small></div>
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
  <v-client-table ref="releasesTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/releases/row.template.jsx'
  import rowChild from 'components/releases/row-child.template.jsx'
  import DeploymentChart from 'components/releases/deployment.chart.vue'

  export default {
    name: 'releases',
    data () {
      return {
        timer: null,
        loaded: false,
        deploymentHeight: 180,
        deploymentData: null,
        stats: {rels: '0', deps: '0'},
        columns: ['ticket_id', 'deployments', 'environments', 'duration', 'begin', 'end'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['ticket_id', 'deployments', 'begin', 'end'],
          dateColumns: ['begin', 'end'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'end', ascending: false},
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
        Vue.$http.get('/release/all')
          .then((response) => {
            if (response != null) {
              this.tableData = response.data.releases
              this.deploymentData = this.fillChart(response.data.chart)
              this.loaded = true
              this.stats = {
                rels: response.data.releases.length.toString(),
                deps: response.data.deployments.toString()
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
