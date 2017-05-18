export default {
    n: function (h, row) {
        return <i class="fa fa-server" aria-hidden="true"></i>
    },
    name: function (h, row) {
        return <a class="text-uppercase" href={'#/vsphere'}>{row.name}</a>
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