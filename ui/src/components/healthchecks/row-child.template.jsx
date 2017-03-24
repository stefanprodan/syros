export default function(h, row) {
   return (
    <div class="row">
        <p>Consul Service ID: {row.service_id}</p>
        <pre>{row.output}</pre>
    </div>
   )
}