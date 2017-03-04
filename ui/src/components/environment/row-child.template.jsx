export default function(h, row) {
   return (
    <div class="row">
        <div class="col-md-6">
            <dl>
                <dt>image</dt>
                <dd>{row.image}</dd>
            </dl>
            <dl>
                <dt>command</dt>
                <dd>{row.command}</dd>
            </dl>
            <dl>
                <dt>labels</dt>
                <dd>  
                    {Object.keys(row.labels).map(function(key){
                        return <p>{key}: {row.labels[key]}</p>              
                    })}
                </dd> 
            </dl>
            <dl>
                <dt>env</dt>
                <dd>  
                    {row.env.map(function(name){
                        return <p>{name}</p>              
                    })}
                </dd> 
            </dl>
        </div>
        <div class="col-md-6">
            <dl>
                <dt>restart_policy</dt>
                <dd>{row.restart_policy}</dd>
            </dl>
            <dl>
                <dt>restart_count</dt>
                <dd>{row.restart_count}</dd>
            </dl>
            <dl>
                <dt>started_at</dt>
                <dd>{row.started_at}</dd>
            </dl>
            <dl>
                <dt>finished_at</dt>
                <dd>{row.finished_at}</dd>
            </dl>
            <dl>
                <dt>exit_code</dt>
                <dd>{row.exit_code}</dd>
            </dl>
            <dl>
                <dt>error</dt>
                <dd>{row.error}</dd>
            </dl>
            <dl>
                <dt>port_bindings</dt>
                <dd>  
                    {Object.keys(row.port_bindings).map(function(key){
                        return <p>{key}: {row.port_bindings[key]}</p>              
                    })}
                </dd> 
            </dl>
            <dl>
                <dt>collected</dt>
                <dd>{row.collected}</dd>
            </dl>
        </div>
    </div>
   )
}