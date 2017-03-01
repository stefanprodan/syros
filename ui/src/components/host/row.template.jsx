export default {
    name: function (h, row) {
        return <a class='' href={'#/container/' + row.id}>{row.name}</a>
    },
    state: function (h, row) {
        if (row.state == 'running'){
            return <span class="alert alert-success text-uppercase">{row.state}</span>
        }
        return <span class="alert alert-danger text-uppercase">{row.state}</span>
    },
    created: function (h, row) {
        return <span>{ moment(row.created).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}