export default {
    n: function (h, row) {
        return <i class="fa fa-sitemap" aria-hidden="true"></i>
    },
    environment: function (h, row) {
        return <a class='' href={'#/environment/' + row.environment}>{row.environment}</a>
    },
    mem_total: function (h, row) {
        return <span>{ parseFloat((row.mem_total / Math.pow(1024, 3))).toFixed(0) }GB</span>
    }
}