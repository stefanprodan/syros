// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import auth from 'components/auth.vue'

// register vue-progressbar
import VueProgressBar from 'vue-progressbar'
Vue.use(VueProgressBar, {
  show: true,
  canSuccess: true,
  color: 'rgb(143, 255, 199)',
  failedColor: 'red',
  height: '2px'
})

// API base URL and auth redirect
import Axios from 'axios'
Axios.defaults.baseURL = process.env.API_LOCATION
Axios.defaults.headers.common.Accept = 'application/json'
Axios.interceptors.response.use(
  response => response,
  (error) => {
    if (error.response != null && error.response.status === 401) {
      auth.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  })
Vue.$http = Axios

// boostrap imports
import jQuery from 'jquery'
window.$ = window.jQuery = jQuery

import moment from 'moment'
window.moment = moment

require('bootstrap')
require('bootstrap/less/bootstrap.less')
require('font-awesome/less/font-awesome.less')
require('./assets/app.scss')

// register vue-table
import {ClientTable} from 'vue-tables-2'
Vue.use(ClientTable, {
  compileTemplates: true,
  highlightMatches: false,
  pagination: {
    dropdown: false,
    chunk: 5
  },
  filterByColumn: false,
  texts: {
    filter: ''
  },
  datepickerOptions: {
    showDropdowns: false
  }
})

// router auth check
router.beforeEach((to, from, next) => {
  if (to.meta.auth && !auth.check()) {
    next({name: 'login', params: {redirect: to.fullPath}})
  } else if (!to.meta.auth && auth.check()) {
    next({name: 'home'})
  } else {
    next()
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
