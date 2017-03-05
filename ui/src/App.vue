<template>
  <div id="app">
    <vue-progress-bar></vue-progress-bar>
    <div class="container-fluid">
      <div class="row">
        <navbar v-if="user.authenticated"></navbar>
        <div class="col-sm-9 col-lg-10">
          <div id="content">
            <flash></flash>
            <router-view></router-view>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import flash from 'components/flash.vue'
  import auth from 'components/auth.vue'
  import navbar from 'components/navbar.vue'

  export default {
    name: 'app',
    data () {
      return {
        user: auth.user
      }
    },
    components: { flash, navbar },
    mounted () {
      //  [App.vue specific] When App.vue is finish loading finish the progress bar
      this.$Progress.finish()
    },
    created () {
      //  [App.vue specific] When App.vue is first loaded start the progress bar
      this.$Progress.start()
      //  hook the progress bar to start before we move router-view
      this.$router.beforeEach((to, from, next) => {
        //  does the page we want to go to have a meta.progress object
        if (to.meta.progress !== undefined) {
          let meta = to.meta.progress
          // parse meta tags
          this.$Progress.parseMeta(meta)
        }
        //  start the progress bar
        this.$Progress.start()
        //  continue to next page
        next()
      })
      //  hook the progress bar to finish after we've finished moving router-view
      this.$router.afterEach((to, from) => {
        //  finish the progress bar
        this.$Progress.finish()
      })
    }
  }
</script>
