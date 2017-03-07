export default function(h, row) {
   return (
    <div class="row">
        <div class="col-md-6">
            {Object.keys(row.config).map(function(key){
                return (
                    <dl>
                        <dt>{key}</dt>
                        <dd>{row.config[key]}</dd>
                    </dl>
                )             
            })}
        </div>
    </div>
   )
}