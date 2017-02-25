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

