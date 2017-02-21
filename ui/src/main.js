// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

import VueProgressBar from 'vue-progressbar'
Vue.use(VueProgressBar, {
  show: true,
  canSuccess: true,
  color: 'rgb(143, 255, 199)',
  failedColor: 'red',
  height: '2px'
})

import Axios from 'axios'
Axios.defaults.baseURL = process.env.API_LOCATION
Axios.defaults.headers.common.Accept = 'application/json'
Axios.interceptors.response.use(
  response => {
    this.$Progress.finish()
    response
  },
  (error) => {
    this.$Progress.fail()
    if (error.response.status === 401) {
      // logout
    }
  })
Vue.$http = Axios

import jQuery from 'jquery'
window.$ = window.jQuery = jQuery

require('bootstrap')
require('bootstrap/less/bootstrap.less')
require('font-awesome/less/font-awesome.less')
require('./assets/app.scss')

import {ClientTable} from 'vue-tables-2'
Vue.use(ClientTable, {
  compileTemplates: false,
  highlightMatches: false,
  pagination: {
    dropdown: false,
    chunk: 5
  },
  filterByColumn: false,
  texts: {
    filter: 'Search:'
  },
  datepickerOptions: {
    showDropdowns: false
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App, ClientTable }
})
