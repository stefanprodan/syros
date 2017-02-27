<script>
import Vue from 'vue'

export default {
  user: {
    authenticated: false
  },
  login (token) {
    this.user.authenticated = true
    localStorage.setItem('token', token)
    Vue.$http.defaults.headers.common.Authorization = `Bearer ${token}`
  },
  logout () {
    this.user.authenticated = true
    localStorage.removeItem('token')
    Vue.$http.defaults.headers.common.Authorization = ''
  },
  check () {
    this.user.authenticated = !!localStorage.getItem('token')
    if (this.user.authenticated) {
      Vue.$http.defaults.headers.common.Authorization = `Bearer ${localStorage.getItem('token')}`
    }

    return this.user.authenticated
  }
}
</script>
