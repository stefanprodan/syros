export default {
    host_name: function (h, row) {
        return <a class='' href={'#/host/' + row.host_id}>{row.host_name}</a>
    },
    status: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 1){
            return <span class="alert alert-danger text-uppercase" title="No signal received for more than one minute ago">removed</span>
        }
        if (row.status == 'passing'){
            return <span class="alert alert-success text-uppercase">{row.status}</span>
        }
        return <span class="alert alert-danger text-uppercase">{row.status}</span>
    },
    since: function (h, row) {
        return <span>{ moment(row.since).fromNow() }</span>
    },
    collected: function (h, row) {
        return <span>{ moment(row.collected).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}