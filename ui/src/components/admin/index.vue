<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li>admin</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.agents }}</h2><small class="text-uppercase">Agents</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.indexers }}</h2><small class="text-uppercase">Indexers</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.rdb }}</h2><small class="text-uppercase">RethinkDB</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.nats }}</h2><small class="text-uppercase">NATS</small></div>
    </div>
  </div>
  <v-client-table ref="servicesTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>  
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/admin/row.template.jsx'
  import rowChild from 'components/admin/row-child.template.jsx'

  export default {
    name: 'admin',
    data () {
      return {
        timer: null,
        stats: {agents: '0', indexers: '0', rdb: '0', nats: '0'},
        columns: ['type', 'status', 'environment', 'hostname', 'collected'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['type', 'environment', 'hostname', 'collected'],
          dateColumns: ['collected'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'type', ascending: true},
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
        Vue.$http.get('/docker/syrosservices')
          .then((response) => {
            if (response != null) {
              this.tableData = response.data
              var agents = 0
              var indexers = 0
              var rdb = 1
              var nats = 1
              for (var i = 0, len = response.data.length; i < len; i++) {
                if (response.data[i].type === 'agent') {
                  agents++
                }
                if (response.data[i].type === 'indexer') {
                  indexers++
                }
              }
              this.stats = {
                agents: agents.toString(),
                indexers: indexers.toString(),
                rdb: rdb.toString(),
                nats: nats.toString()
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
