export default {
    status: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 1){
            return <span class="alert alert-danger text-uppercase" title="No signal received for more than one minute ago">removed</span>
        }
        if (row.status == 'passing'){
            return <span class="alert alert-success text-uppercase">{row.status}</span>
        }
        return <span class="alert alert-danger text-uppercase">{row.status}</span>
    },
    duration: function (h, row) {
        var begin = moment(row.begin)
        var end = moment(row.end)
        return <span>{ moment(row.end).diff(row.begin, 'minutes')} minutes</span>
    },
    begin: function (h, row) {
        return <span>{ moment(row.begin).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    },
    end: function (h, row) {
        return <span>{ moment(row.end).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}