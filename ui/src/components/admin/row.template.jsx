export default {
    status: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 1){
            return <span class="alert alert-danger text-uppercase" title="No signal received for more than one minute ago">down</span>
        }
        return <span class="alert alert-success text-uppercase" title="Signal received less than one minute ago">up</span>
    },
    collected: function (h, row) {
        return <span>{ moment(row.collected).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}