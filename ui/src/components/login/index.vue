<template>
  <div class="row login">
    <div class="col-md-6 col-md-offset-3">
      <h1>Syros</h1>
      <div class="alert alert-danger" v-if="error">
        <p>{{ error }}</p>
      </div>
      <form @submit.prevent="login(user)">
        <div class="form-group">
          <div class="input-group">
            <div class="input-group-addon">
              <i class="fa fa-envelope fa-fw"></i>
            </div>
            <input
              v-model="user.email"
              type="input"
              placeholder="Email"
              class="form-control"
            >
          </div>
        </div>
        <div class="form-group">
          <div class="input-group">
            <div class="input-group-addon">
              <i class="fa fa-lock fa-fw"></i>
            </div>
            <input
              v-model="user.password"
              type="password"
              placeholder="Password"
              class="form-control"
            >
          </div>
        </div>
        <div class="form-group">
          <button class="btn btn-primary">
            Login
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
  import auth from 'components/auth.vue'

  export default {
    name: 'login',
    data () {
      return {
        error: null,
        user: {
          email: null,
          password: null
        }
      }
    },
    methods: {
      login (user) {
        if (!user.email || !user.password) {
          this.error = 'Invalid Username or Password'
        } else {
          auth.login(user.email)
          this.$router.push('/')
        }
      }
    }
  }

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1,
  h2 {
    font-weight: normal;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
  }
  
  li {
    display: inline-block;
    margin: 0 10px;
  }
  
  a {
    color: #42b983;
  }
</style>
