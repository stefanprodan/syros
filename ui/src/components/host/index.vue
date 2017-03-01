<template>
<div>
  <div id="main-breadcrumb">
    <ol class="breadcrumb">
      <li><a class="text-uppercase" href="/#/home">Home</a></li>
      <li><a class="text-uppercase" href="/#/hosts">hosts</a></li>
      <li>{{ stats.name }}</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.containers }}</h2><small class="text-uppercase">Containers</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.images }}</h2><small class="text-uppercase">Images</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.cpus }}</h2><small class="text-uppercase">vCPUs</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.ram }}</h2><small class="text-uppercase">Memory</small></div>
    </div>
  </div>

  <v-client-table ref="containersTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>
  
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/host/row.template.jsx'
  import rowChild from 'components/host/row-child.template.jsx'

  export default {
    name: 'host',
    data () {
      return {
        timer: null,
        id: null,
        stats: {name: '', containers: '0', images: '0', cpus: '0', ram: '0 MB'},
        columns: ['name', 'state', 'status', 'network_mode', 'port', 'created'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['name', 'state', 'status', 'network_mode', 'port', 'created'],
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
    methods: {
      loadData () {
        this.$Progress.start()
        Vue.$http.get(`/docker/hosts/${this.id}`)
          .then((response) => {
            if (response != null) {
              this.tableData = response.data.Containers
              this.stats = {
                name: response.data.Host.name,
                containers: response.data.Host.containers_running.toString(),
                images: response.data.Host.images.toString(),
                cpus: response.data.Host.ncpu.toString(),
                ram: parseInt(parseFloat((response.data.Host.mem_total / Math.pow(1024, 3))).toFixed(0)).toString() + 'GB'
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
      this.id = this.$route.params.id
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
