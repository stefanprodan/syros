<template>
  <transition name="fade">
    <div
      v-if="visible"
      v-bind:class="type"
      role="alert"
      v-text="message"
    >
    </div>
  </transition>
</template>

<style scoped>
  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s
  }
  .fade-enter, .fade-leave-to {
    opacity: 0
  }
</style>

<script>
import bus from 'components/bus.vue'

export default {
  name: 'flash-alert',
  data () {
    return {
      type: '',
      message: '',
      visible: false
    }
  },
  created () {
    bus.$on('flashMessage', data => this.setData(data))
  },
  methods: {
    setData (data) {
      this.type = `flash-alert alert alert-${data.type}`
      this.message = data.message
      this.visible = true
    },
    setFadeOut () {
      setTimeout(() => (
        this.visible = false
      ), 5000)
    }
  },
  watch: {
    visible: 'setFadeOut'
  }
}
</script>
