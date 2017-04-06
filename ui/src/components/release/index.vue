<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><router-link class="text-uppercase" :to="{ name: 'home'}">home</router-link></li>
      <li><router-link class="text-uppercase" :to="{ name: 'releases'}">releases</router-link></li>
      <li>{{ stats.name }}</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-6 text-center">
        <h2>{{ stats.deps }}</h2><small class="text-uppercase">Deployments</small></div>
      <div class="col-md-6 text-center">
        <h2>{{ stats.name }}</h2><small class="text-uppercase">Release</small></div>
    </div>
  </div>

  <v-client-table ref="releaseTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>  
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import rowTemplate from 'components/release/row.template.jsx'
  import rowChild from 'components/release/row-child.template.jsx'

  export default {
    name: 'release',
    data () {
      return {
        timer: null,
        id: null,
        loaded: false,
        envHeight: 180,
        stats: {deps: '0', name: ''},
        columns: ['service_name', 'status', 'host_name', 'environment', 'timestamp'],
        tableData: [],
        options: {
          skin: 'table-hover',
          sortable: ['service_name', 'status', 'host_name', 'environment', 'timestamp'],
          dateColumns: ['timestamp'],
          toMomentFormat: 'YYYY-MM-DDTHH:mm:ssZ',
          uniqueKey: 'id',
          orderBy: {column: 'timestamp', ascending: true},
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
        Vue.$http.get(`/release/${this.id}`)
          .then((response) => {
            if (response != null) {
              this.tableData = response.data
              this.stats = {
                name: response.data[0].ticket_id,
                deps: response.data.length.toString()
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
