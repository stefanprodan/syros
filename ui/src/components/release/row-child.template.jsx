export default function(h, row) {
    if (row.log.length < 1){
        row.log = row.host_name + ':'+ row.service_name
    }
   return (
    <div class="row">
        <pre>{row.log}</pre>
    </div>
   )
}