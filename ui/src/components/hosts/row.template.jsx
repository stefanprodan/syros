export default {
    name: function (h, row) {
        return <a class='fa fa-edit' href={'#/' + row.id + '/name'}>{row.name}</a>
    },
    age: function (h, row) {
        return <a class='fa fa-edit' href={'#/' + row.id + '/age'}>{row.age}</a>
    }
}