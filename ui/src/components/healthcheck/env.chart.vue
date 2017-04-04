<script>
import {Bar, mixins} from 'vue-chartjs'
export default Bar.extend({
  mixins: [mixins.reactiveProp],
  props: ['chartData'],
  mounted () {
    this.renderChart(this.chartData, {
      responsive: true,
      maintainAspectRatio: false,
      legend: {
        display: false
      },
      tooltips: {
        mode: 'index',
        intersect: false
      },
      scales: {
        yAxes: [{ stacked: true, ticks: { beginAtZero: true } }],
        xAxes: [{ stacked: true, ticks: { beginAtZero: true } }]
      }
    })
  }
})
</script>
