export default function(h, row) {
    if (row.output.length < 1){
        row.output = 'TIMEOUT'
    }
   return (
    <div class="row">
        <pre>{row.output}</pre>
    </div>
   )
}