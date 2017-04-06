export default {
    duration: function (h, row) {
        var end = moment(row.end)
        var duration = end.diff(row.begin, 'minutes') + ' minutes'
        if (end.diff(row.begin, 'minutes') < 2){
            duration = end.diff(row.begin, 'seconds') + ' seconds'
        }
        if (end.diff(row.begin, 'minutes') > 120){
            duration = end.diff(row.begin, 'hours') + ' hours'
        }
        if (end.diff(row.begin, 'h') > 48){
            duration = end.diff(row.begin, 'days') + ' days'
        }

        return <span>{ duration }</span>
    },
    begin: function (h, row) {
        return <span>{ moment(row.begin).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    },
    end: function (h, row) {
        return <span>{ moment(row.end).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
}