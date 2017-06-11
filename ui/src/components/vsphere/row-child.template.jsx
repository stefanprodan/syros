export default function(h, row) {
   return (
    <div class="row">
        <div class="col-md-6">
            <dl>
                <dt>cluster</dt>
                <dd>{row.cluster}</dd>
            </dl>
            <dl>
                <dt>host_id</dt>
                <dd>{row.host_id}</dd>
            </dl>
            <dl>
                <dt>host_name</dt>
                <dd>{row.host_name}</dd>
            </dl>
            <dl>
                <dt>datastore_id</dt>
                <dd>{row.datastore_id}</dd>
            </dl>
            <dl>
                <dt>datastore_name</dt>
                <dd>{row.datastore_name}</dd>
            </dl>
        </div>
        <div class="col-md-6">
            <dl>
                <dt>ID</dt>
                <dd>{row.id}</dd>
            </dl>
            <dl>
                <dt>environment</dt>
                <dd>{row.environment}</dd>
            </dl>
            <dl>
                <dt>IP</dt>
                <dd>{row.ip}</dd>
            </dl>
            <dl>
                <dt>boot_time</dt>
                <dd>{(row.boot_time) ? moment(row.boot_time).format('YYYY-MM-DD HH:mm:ss Z') : 'unknown'}</dd>
            </dl>
            <dl>
                <dt>collected</dt>
                <dd>{moment(row.collected).format('YYYY-MM-DD HH:mm:ss Z')}</dd>
            </dl>
        </div>
    </div>
   )
}