<script>
import {Pie, mixins} from 'vue-chartjs'
const { reactiveData } = mixins
export default Pie.extend({
  mixins: [reactiveData],
  data () {
    return {
      chartData: null
    }
  },
  mounted () {
    this.renderChart(this.chartData, {responsive: true, maintainAspectRatio: false})
  },
  methods: {
    fillData (data) {
      this.chartData = data
    }
  }
})
</script>
