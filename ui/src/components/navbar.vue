<template>
<div class="col-sm-3 col-lg-2">
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
            <a><i class="fa fa-home"></i> Home</a>
        </router-link> 
        <router-link :to="{ name: 'healthchecks' }" active-class="active" tag="li">
            <a><i class="fa fa-heartbeat"></i> Health</a>
        </router-link>
        <router-link :to="{ name: 'releases' }" active-class="active" tag="li">
            <a><i class="fa fa-play-circle-o"></i> Releases</a>
        </router-link>
        <router-link :to="{ name: 'vsphere' }" active-class="active" tag="li">
            <a><i class="fa fa-sitemap"></i> vSphere</a>
        </router-link>
        <router-link :to="{ name: 'hosts' }" active-class="active" tag="li">
            <a><i class="fa fa-h-square"></i> Docker Hosts</a>
        </router-link> 
        </ul>
        <ul v-if="environments" class="nav navbar-nav navbar-indent">
            <li class="dropdown-header">Environments</li>
            <router-link v-for="env in environments" active-class="active" tag="li" :to="{ name: 'environment', params: { id: env }}">
            <a>{{ env }}</a>
            </router-link>
        </ul>
        <ul class="nav navbar-nav">   
            <router-link :to="{ name: 'admin' }" active-class="active" tag="li">
                <a><i class="fa fa-circle-o-notch"></i> Admin</a>
            </router-link>  
            <li><a href="#" @click="logout()"><i class="fa fa-sign-out"></i> Logout</a></li>
        </ul>
        <p class="navbar-text">
            <a href="https://github.com/stefanprodan/syros">Syros v0.7</a> open-source software by <a href="http://www.stefanprodan.com">Stefan</a>
        </p>
        </div>  
    </nav>
</div>
</template>

<script>
  import Vue from 'vue'
  import bus from 'components/bus.vue'
  import auth from 'components/auth.vue'

  export default {
    name: 'navbar',
    data () {
      return {
        environments: []
      }
    },
    methods: {
      loadData () {
        Vue.$http.get('/home/environments')
          .then((response) => {
            if (response != null) {
              this.environments = response.data
            }
          })
          .catch((error) => {
            if (error.message === 'Network Error') {
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
    mounted () {
      this.refreshData()
    }
  }
</script>
