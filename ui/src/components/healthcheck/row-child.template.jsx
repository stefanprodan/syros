export default function(h, row) {
    if (row.output.length < 1){
        row.output = 'TIMEOUT'
    }
   return (
    <div class="row">
        <p>Check ID: {row.service_id}</p>
        <pre>{row.output}</pre>
    </div>
   )
}