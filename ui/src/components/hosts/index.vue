<template>
<div>
  <div id="main-breadcrumb">
    <ol class="breadcrumb">
      <li><a class="text-uppercase" href="/#/home">Home</a></li>
      <li><a class="text-uppercase" href="/#/hosts">hosts</a></li>
      <li>prep-assetengine-2</li>
    </ol>
  </div>
  <div class="stats">
    <div class="row">
      <div class="col-md-3 text-center">
        <h2>{{ stats.hosts }}</h2><small class="text-uppercase">Hosts</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.containers }}</h2><small class="text-uppercase">Containers</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.cpus }}</h2><small class="text-uppercase">vCPUs</small></div>
      <div class="col-md-3 text-center">
        <h2>{{ stats.ram }}</h2><small class="text-uppercase">Memory</small></div>
    </div>
  </div>

  <v-client-table ref="hostsTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>

   <v-client-table ref="hostsTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>

  
</div>
</template>

<script>
  import rowTemplate from 'components/hosts/row.template.jsx'
  import rowChild from 'components/hosts/row-child.template.jsx'

  export default {
    name: 'hosts',
    data () {
      return {
        columns: ['name', 'status', 'age', 'date', 'edit'],
        stats: {hosts: '12', containers: '34', cpus: '24', ram: '12064 MB'},
        tableData: [
            {id: '1', name: 'John', status: 'up', age: '33', date: '21/02/2015'},
            {id: '2', name: 'Jane', status: 'up', age: '24', date: '18/02/2015'},
            {id: '3', name: 'Susan', status: 'up', age: '16', date: '21/03/2016'},
            {id: '4', name: 'Chris', status: 'up', age: '55', date: '11/05/2017'},
            {id: '5', name: 'John', status: 'down', age: '20', date: '11/02/2015'},
            {id: '6', name: 'Jane', status: 'up', age: '24', date: '12/02/2015'},
            {id: '7', name: 'Susan', status: 'up', age: '16', date: '13/02/2015'},
            {id: '8', name: 'Chris', status: 'up', age: '55', date: '14/02/2015'},
            {id: '9', name: 'John', status: 'up', age: '20', date: '15/02/2015'},
            {id: '10', name: 'Jane', status: 'up', age: '24', date: '16/02/2015'},
            {id: '11', name: 'Susan', status: 'up', age: '16', date: '17/02/2015'},
            {id: '12', name: 'Chris', status: 'up', age: '55', date: '18/02/2015'},
            {id: '13', name: 'John', status: 'up', age: '20', date: '19/02/2015'},
            {id: '14', name: 'Jane', status: 'up', age: '24', date: '20/02/2015'},
            {id: '15', name: 'Susan', status: 'up', age: '16', date: '23/02/2015'},
            {id: '16', name: 'Chris', status: 'up', age: '55', date: '24/02/2015'},
            {id: '17', name: 'Dan', status: 'up', age: '40', date: '26/02/2018'}
        ],
        options: {
          skin: 'table-hover',
          sortable: ['name', 'status', 'age', 'date'],
          dateColumns: ['date'],
          toMomentFormat: 'DD/MM/YYYY', // YYYY-MM-DDTHH:mm:ssZ
          uniqueKey: 'id',
          orderBy: {column: 'date', ascending: true},
          perPage: 3,
          perPageValues: [10, 20, 30, 50],
          childRow: rowChild,
          templates: rowTemplate
        }
      }
    },
    methods: {
      toggleChildRow (id) {
        this.$refs.hostsTabel.toggleChildRow(id)
      },
      hello () {
        alert('This is the message.')
      }
    },
    created: function () {
      console.log('Created: ' + this.$options.name)
    },
    mounted: function () {
      console.log('Mounted: ' + this.$options.name)

      this.$on('toggle', function (id) {
        console.log(id)
        this.$refs.hostsTabel.toggleChildRow(id)
      })

      // setTimeout(
      //   () => {
      //     this.tableData = [
      //       {id: '1', name: 'John', status: 'up', age: '33', date: '21/02/2015'},
      //       {id: '2', name: 'Jane', status: 'up', age: '24', date: '18/02/2015'},
      //       {id: '3', name: 'Susan', status: 'up', age: '16', date: '21/03/2016'},
      //       {id: '4', name: 'Chris', status: 'up', age: '55', date: '11/05/2017'},
      //       {id: '5', name: 'John', status: 'down', age: '20', date: '11/02/2015'},
      //       {id: '6', name: 'Jane', status: 'up', age: '24', date: '12/02/2015'},
      //       {id: '7', name: 'Susan', status: 'up', age: '16', date: '13/02/2015'},
      //       {id: '8', name: 'Chris', status: 'up', age: '55', date: '14/02/2015'},
      //       {id: '9', name: 'John', status: 'up', age: '20', date: '15/02/2015'},
      //       {id: '10', name: 'Jane', status: 'up', age: '24', date: '16/02/2015'},
      //       {id: '11', name: 'Susan', status: 'up', age: '16', date: '17/02/2015'},
      //       {id: '12', name: 'Chris', status: 'up', age: '55', date: '18/02/2015'},
      //       {id: '13', name: 'John', status: 'up', age: '20', date: '19/02/2015'},
      //       {id: '14', name: 'Jane', status: 'up', age: '24', date: '20/02/2015'},
      //       {id: '15', name: 'Susan', status: 'up', age: '16', date: '23/02/2015'},
      //       {id: '16', name: 'Chris', status: 'up', age: '55', date: '24/02/2015'},
      //       {id: '17', name: 'Dan', status: 'up', age: '40', date: '26/02/2018'}
      //     ]
      //   },
      //   500
      // )
    },
    destroyed: function () {
      console.log('Destroyed: ' + this.$options.name)
    }
}

</script>
