export default {
    n: function (h, row) {
        return <i class="fa fa-usb" aria-hidden="true"></i>
    },
    power_state: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 5){
            return <span class="alert alert-warning text-uppercase" title="No signal received for more than one minute ago">unknown</span>
        }
        return <span class="alert alert-success text-uppercase">POWEREDON</span>
    },
    capacity: function (h, row) {
        return <span>{ parseFloat((row.capacity / Math.pow(1024, 3))).toFixed(0) }GB</span>
    },
    free: function (h, row) {
        return <span>{ parseFloat((row.free / Math.pow(1024, 3))).toFixed(0) }GB</span>
    }
}