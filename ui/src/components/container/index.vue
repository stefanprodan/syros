<template>
<div>
  <div>
    <ol class="breadcrumb">
      <li><a class="text-uppercase" href="/#/home">Home</a></li>
      <li><router-link class="text-uppercase" :to="{ name: 'hosts', params: { id: stats.host_id }}">{{stats.host}}</router-link></li>
      <li>{{ stats.name }}</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.state }}</h2><small class="text-uppercase">state</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.status }}</h2><small class="text-uppercase">status</small></div>
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
  import rowTemplate from 'components/container/row.template.jsx'

  export default {
    name: 'host',
    data () {
      return {
        timer: null,
        id: null,
        stats: {name: '', host: '', host_id: '', state: '', status: '', cpus: '0', ram: '0 MB'},
        columns: ['n', 'property', 'value'],
        tableData: [],
        options: {
          skin: 'table-hover',
          uniqueKey: 'n',
          perPage: 10,
          perPageValues: [10, 20, 30, 50],
          templates: rowTemplate
        }
      }
    },
    methods: {
      parseProps (item) {
        var n = 1
        var tabelProps = [
          {n: n++, icon: 'fa fa-info', property: 'name', value: item.name},
          {n: n++, icon: 'fa fa-heartbeat', property: 'state', value: item.state},
          {n: n++, icon: 'fa fa-info-circle', property: 'status', value: item.status},
          {n: n++, icon: 'fa fa-inbox', property: 'image', value: item.image},
          {n: n++, icon: 'fa fa-terminal', property: 'command', value: item.command},
          {n: n++, icon: 'fa fa-clock-o', property: 'created', value: item.created},
          {n: n++, icon: 'fa fa-clock-o', property: 'collected', value: item.collected},
          {n: n++, icon: 'fa fa-clock-o', property: 'started_at', value: item.started_at},
          {n: n++, icon: 'fa fa-clock-o', property: 'finished_at', value: item.finished_at},
          {n: n++, icon: 'fa fa-repeat', property: 'restart_policy', value: item.restart_policy},
          {n: n++, icon: 'fa fa-refresh', property: 'restart_count', value: item.restart_count},
          {n: n++, icon: 'fa fa-info-circle', property: 'exit_code', value: item.exit_code},
          {n: n++, icon: 'fa fa-exclamation-triangle', property: 'error', value: item.error},
          {n: n++, icon: 'fa fa-sitemap', property: 'network_mode', value: item.network_mode}
        ]

        Object.keys(item.port_bindings).map(function (key) {
          n++
          tabelProps.push({n: n, icon: 'fa fa-sitemap', property: 'port_binding: ' + key, value: item.port_bindings[key]})
        })
        Object.keys(item.labels).map(function (key) {
          n++
          tabelProps.push({n: n, icon: 'fa fa-tag', property: 'label: ' + key, value: item.labels[key]})
        })
        item.env.map(function (name) {
          n++
          tabelProps.push({n: n, icon: 'fa fa-clipboard', property: 'env: ' + name.split('=')[0], value: name.split('=')[1]})
        })
        this.tableData = tabelProps
      },
      loadData () {
        this.$Progress.start()
        Vue.$http.get(`/docker/containers/${this.id}`)
          .then((response) => {
            if (response != null) {
              this.parseProps(response.data.containers[0])
              this.stats = {
                host: response.data.host.name,
                host_id: response.data.host.id,
                name: response.data.containers[0].name,
                state: response.data.containers[0].state,
                status: response.data.containers[0].status,
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
                icon: 'warning',
                message: 'Network error! Could not connect to the server'
              })
            } else {
              bus.$emit('flashMessage', {
                icon: 'warning',
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
