<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li><router-link class="text-uppercase" :to="{ name: 'healthchecks'}">health</router-link></li>
      <li>{{ stats.name }}</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2 class="critical">{{ stats.criticalCount }}</h2><small class="text-uppercase">Critical</small></div>
      <div class="col-md-3 text-center">
        <h2 class="critical">{{ stats.criticalDuration }}</h2><small class="text-uppercase">Critical last 30d</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.passingCount }}</h2><small class="text-uppercase">Passing</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.passingDuration }}</h2><small class="text-uppercase">Passing last 30d</small></div>
    </div>
  </div>

  <v-client-table ref="healthchecksTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>  
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/healthcheck/row.template.jsx'
  import rowChild from 'components/healthcheck/row-child.template.jsx'

  export default {
    name: 'healthchecks',
    data () {
      return {
        timer: null,
        id: null,
        loaded: false,
        envHeight: 180,
        stats: {criticalCount: '0', criticalDuration: '0s', passingCount: '0', passingDuration: '0s', name: ''},
        columns: ['status', 'duration', 'begin', 'end'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['begin'],
          dateColumns: ['begin', 'end'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'begin', ascending: false},
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          childRow: rowChild,
          templates: rowTemplate
        }
      }
    },
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get(`/docker/healthchecks/${this.id}`)
          .then((response) => {
            if (response != null) {
              this.tableData = response.data.checks
              var statsCriticalCount = 0
              var statsCriticalDuration = '0s'
              var statsPassingCount = 0
              var statsPassingDuration = '0s'

              for (var i = 0, len = response.data.stats.length; i < len; i++) {
                if (response.data.stats[i].status === 'passing') {
                  statsPassingCount = response.data.stats[i].count
                  statsPassingDuration = window.moment.duration(response.data.stats[i].duration, 'seconds').humanize()
                } else {
                  statsCriticalCount = response.data.stats[i].count
                  statsCriticalDuration = window.moment.duration(response.data.stats[i].duration, 'seconds').humanize()
                }
              }
              this.stats = {
                name: response.data.checks[0].service_id,
                criticalCount: statsCriticalCount.toString(),
                criticalDuration: statsCriticalDuration,
                passingCount: statsPassingCount.toString(),
                passingDuration: statsPassingDuration
              }

              this.loaded = true

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
      this.id = this.$route.params.id
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
