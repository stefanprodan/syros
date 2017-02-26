export default {
    name: function (h, row) {
        return <a class='' href={'#/host/' + row.id}>{row.name}</a>
    },
    status: function (h, row) {
        if (moment().diff(row.collected, 'minutes') > 1){
            return <span class="alert alert-danger text-uppercase" title="No signal received for more than one minute">down</span>
        }
        return <span class="alert alert-success text-uppercase" title="Signal received less than one minute">up</span>
    },
    mem_total: function (h, row) {
        return <span>{ parseFloat((row.mem_total / Math.pow(1024, 3))).toFixed(0) }GB</span>
    },
    system_time: function (h, row) {
        return <span>{ moment(row.system_time).format('YYYY-MM-DD HH:mm:ss Z') }</span>
    }
    /*edit: function (h, row) {
        return <a href="#" on-click={ () => this.$emit('toggle',row.id) }><i class="fa fa-pencil"></i></a>
    }
    age: function (h, row) {
        console.log(this)
        return (
        <button on-click={ () => this.hello() } class="btn btn-primary btn-sm">
            <i class="glyphicon glyphicon-edit"></i>
        </button>
        )
    }*/
    // edit: function (h, row) {
    //     return <button on-click={ () => this.$refs.hostsTabel.toggleChildRow(row.id) } class="btn btn-primary btn-sm"><i class="glyphicon glyphicon-edit"></i></button>
    // }
}