import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/home',
      name: 'home',
      component: require('components/home/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/releases',
      name: 'releases',
      component: require('components/releases/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/release/:id',
      name: 'release',
      component: require('components/release/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/healthchecks',
      name: 'healthchecks',
      component: require('components/healthchecks/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/healthcheck/:id',
      name: 'healthcheck',
      component: require('components/healthcheck/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/hosts',
      name: 'hosts',
      component: require('components/hosts/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/vsphere',
      name: 'vsphere',
      component: require('components/vsphere/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/host/:id',
      name: 'host',
      component: require('components/host/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/container/:id',
      name: 'container',
      component: require('components/container/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/environment/:id',
      name: 'environment',
      component: require('components/environment/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/admin',
      name: 'admin',
      component: require('components/admin/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/login/:redirect?',
      name: 'login',
      component: require('components/login/index.vue'),
      meta: {
        auth: false
      }
    },
    {
      path: '/',
      redirect: '/home'
    },
    {
      path: '/*',
      redirect: '/home'
    }
  ]
})
