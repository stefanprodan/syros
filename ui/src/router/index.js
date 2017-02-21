import Vue from 'vue'
import Router from 'vue-router'
import Hello from 'components/Hello'
import Host from 'components/Host'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/home',
      name: 'Hello',
      component: Hello
    },
    {
      path: '/hosts',
      name: 'hosts.index',
      component: require('components/hosts/index.vue')
    },
    {
      path: '/host',
      name: 'Host',
      component: Host
    }
  ]
})
