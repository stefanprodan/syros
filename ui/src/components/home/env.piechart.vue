<script>
import {Doughnut, mixins} from 'vue-chartjs'
export default Doughnut.extend({
  mixins: [mixins.reactiveProp],
  props: ['chartData'],
  mounted () {
    this.renderChart(this.chartData, {responsive: true, maintainAspectRatio: true})
  }
})
</script>
