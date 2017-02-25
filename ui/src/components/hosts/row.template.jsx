export default {
    name: function (h, row) {
        return <a class='' href={'#/' + row.id + '/name'}>{row.name}</a>
    },
    status: function (h, row) {
        if (row.status == 'down'){
            return <span class="alert alert-danger text-uppercase">{row.status}</span>
        }
        return <span class="alert alert-success text-uppercase">{row.status}</span>
    },
    age: function (h, row) {
        return (
            <a class='' href={'#/' + row.id + '/name'}>{row.age}</a>
        )
    },
    edit: function (h, row) {
        return <a href="#" on-click={ () => this.$emit('toggle',row.id) }><i class="fa fa-pencil"></i></a>
    }
    /*age: function (h, row) {
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