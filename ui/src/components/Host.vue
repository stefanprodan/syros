<template>
  <v-select :value.sync="selected" :options="options"></v-select>
</template>

<script>
  import vSelect from 'vue-select'

  export default {
    name: 'host',
    data () {
      return {
        selected: null,
        options: ['foo', 'bar', 'baz']
      }
    },
    components: { vSelect }
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
</style>
