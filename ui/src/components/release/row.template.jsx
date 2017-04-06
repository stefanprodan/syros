export default {
    status: function (h, row) {
        return <span class="alert alert-success text-uppercase">{row.status}</span>
    },
    timestamp: function (h, row) {
        return <span>{ moment(row.end).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}