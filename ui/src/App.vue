<template>
  <div id="app">
    <vue-progress-bar></vue-progress-bar>
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-3 col-lg-2">
          <nav class="navbar navbar-default navbar-fixed-side navbar-inverse">
            
              <div class="navbar-header">
                <button class="navbar-toggle" data-target=".navbar-collapse" data-toggle="collapse">
                  <span class="sr-only">Toggle navigation</span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                </button>
                <router-link class="navbar-brand" :to="{ name: 'Hello' }">Syros</router-link>
              </div>
              <div class="collapse navbar-collapse">
                <ul class="nav navbar-nav">
                <router-link
                  :to="{ name: 'Hello' }"
                  active-class="active"
                  tag="li"
                >
                  <a>
                    Home
                  </a>
                </router-link>
                <router-link
                  :to="{ name: 'hosts.index' }"
                  active-class="active"
                  tag="li"
                >
                  <a>
                    Hosts
                  </a>
                </router-link>                
                </ul>
                <ul class="nav navbar-nav navbar-indent">
                  <li class="dropdown-header">Environments</li>
                  <li><a href="#">INT</a></li>
                  <li><a href="#">STG</a></li>
                  <li><a href="#">PREP</a></li>
                  <li><a href="#">PROD</a></li>
                </ul>
                <ul class="nav navbar-nav">
                  <li><a href="#"><i class="fa fa-sign-out"></i> Logout</a></li>
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
            <router-view></router-view>
          </div>
        </div>
      </div>
    </div>


  </div>
</template>

<script>
  export default {
    name: 'app',
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
