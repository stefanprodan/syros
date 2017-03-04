<template>
  <div id="app">
    <vue-progress-bar></vue-progress-bar>
    <div class="container-fluid">
      <div class="row">
        <div v-if="user.authenticated" class="col-sm-3 col-lg-2">
          <nav class="navbar navbar-default navbar-fixed-side navbar-inverse">
            
              <div class="navbar-header">
                <button class="navbar-toggle" data-target=".navbar-collapse" data-toggle="collapse">
                  <span class="sr-only">Toggle navigation</span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                </button>
                <router-link class="navbar-brand" :to="{ name: 'home' }">Syros</router-link>
              </div>
              <div class="collapse navbar-collapse">
                <ul class="nav navbar-nav">
                <router-link :to="{ name: 'home' }" active-class="active" tag="li">
                  <a>Home</a>
                </router-link>
                <router-link :to="{ name: 'hosts' }" active-class="active" tag="li">
                  <a>Hosts</a>
                </router-link>                
                </ul>
                <ul v-if="environments" class="nav navbar-nav navbar-indent">
                  <li class="dropdown-header">Environments</li>
                  <router-link v-for="env in environments" active-class="active" tag="li" :to="{ name: 'environment', params: { id: env }}">
                    <a>{{ env }}</a>
                  </router-link>
                </ul>
                <ul class="nav navbar-nav">
                  <li><a href="#" v-if="user.authenticated" @click="logout()"><i class="fa fa-sign-out"></i> Logout</a></li>
                </ul>
                <p class="navbar-text">
                  Made by
                  <a href="http://www.stefanprodan.com">Stefan</a>
                </p>
              </div>
            
          </nav>
        </div>
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
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import flash from 'components/flash.vue'
  import auth from 'components/auth.vue'

  export default {
    name: 'app',
    data () {
      return {
        user: auth.user,
        environments: []
      }
    },
    methods: {
      loadData () {
        Vue.$http.get('/docker/environments')
          .then((response) => {
            if (response != null) {
              this.environments = response.data
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
          })
      },
      refreshData () {
        this.loadData()
        console.log('Refresh data: ' + this.$options.name)
        // enqueue new call after 30 seconds
        if (this.timer) clearTimeout(this.timer)
        this.timer = setTimeout(this.refreshData, 30000)
      },
      logout () {
        auth.logout()
        this.$router.push({
          name: 'login'
        })
      }
    },
    components: { flash },
    mounted () {
      this.refreshData()
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
