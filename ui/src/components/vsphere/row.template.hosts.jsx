export default {
    n: function (h, row) {
        return <i class="fa fa-server" aria-hidden="true"></i>
    },
    power_state: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 5){
            return <span class="alert alert-warning text-uppercase" title="No signal received for more than one minute ago">unknown</span>
        }
        if (row.power_state == 'poweredOn'){
            return <span class="alert alert-success text-uppercase">{row.power_state}</span>
        }
        return <span class="alert alert-danger text-uppercase">{row.power_state}</span>
    },
    memory: function (h, row) {
        return <span>{ parseInt(parseFloat((row.memory / Math.pow(1024, 3))).toFixed(0)) }GB</span>
    },
    boot_time: function (h, row) {
        if (!row.boot_time){
            return <span>unknown</span>
        }
        return <span>{ moment(row.boot_time).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}