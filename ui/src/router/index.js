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
      path: '/hosts',
      name: 'hosts',
      component: require('components/hosts/index.vue'),
      meta: {
        auth: true
      }
    },
    {
      path: '/login',
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
