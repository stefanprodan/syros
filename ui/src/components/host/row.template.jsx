export default {
    name: function (h, row) {
        return <a class='' href={'#/container/' + row.id}>{row.name}</a>
    },
    state: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 1){
            return <span class="alert alert-warning text-uppercase" title="No signal received for more than one minute ago">removed</span>
        }
        if (row.state == 'running'){
            return <span class="alert alert-success text-uppercase">{row.state}</span>
        }
        return <span class="alert alert-danger text-uppercase">{row.state}</span>
    },
    created: function (h, row) {
        return <span>{ moment(row.created).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}