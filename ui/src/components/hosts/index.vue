<template>
  <v-client-table ref="hostsTabel" :data="tableData" :columns="columns" :options="options"></v-client-table>
</template>

<script>
  import rowTemplate from 'components/hosts/row.template.jsx'
  import rowChildTemplate from 'components/hosts/row-child.template.jsx'

  export default {
    name: 'hosts-index',
    data () {
      return {
        columns: ['id', 'name', 'age'],
        tableData: [
          {id: '1', name: 'John', age: '33'},
          {id: '2', name: 'Jane', age: '24'},
          {id: '3', name: 'Susan', age: '16'},
          {id: '4', name: 'Chris', age: '55'},
          {id: '5', name: 'John', age: '20'},
          {id: '6', name: 'Jane', age: '24'},
          {id: '7', name: 'Susan', age: '16'},
          {id: '8', name: 'Chris', age: '55'},
          {id: '9', name: 'John', age: '20'},
          {id: '10', name: 'Jane', age: '24'},
          {id: '11', name: 'Susan', age: '16'},
          {id: '12', name: 'Chris', age: '55'},
          {id: '13', name: 'John', age: '20'},
          {id: '14', name: 'Jane', age: '24'},
          {id: '15', name: 'Susan', age: '16'},
          {id: '16', name: 'Chris', age: '55'},
          {id: '17', name: 'Dan', age: '40'}
        ],
        options: {
          sortable: ['id', 'name', 'age'],
          templates: rowTemplate,
          childRow: rowChildTemplate
        }
      }
    },
    methods: {
      formatLink (value) {
        return '<a href="#">' + value + '</a>'
      }
    },
    created: function () {
      console.log('Created: ' + this.$options.name)
    },
    mounted: function () {
      console.log('Mounted: ' + this.$options.name)
      this.$refs.hostsTabel.toggleChildRow(1)
    },
    destroyed: function () {
      console.log('Destroyed: ' + this.$options.name)
    }
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  #vue-tables-options {
    text-align: center;
  }
  
  table {
    text-align: left;
  }
  
  h2 {
    margin-bottom: 10px;
  }
  
  thead tr:nth-child(2) th {
    font-weight: normal;
  }
  
  .VueTables__sortable {
    cursor: pointer;
  }

.VueTables__child-row-toggler--closed::before {
  width:16px;
  height:16px;
  line-height: 16px;
  display: block;
  text-align: center;
  content: "+";
}

.VueTables__child-row-toggler--open::before  {
  width:16px;
  height:16px;
  line-height: 16px;
  display: block;
  text-align: center;
  content: "-";
}
</style>
