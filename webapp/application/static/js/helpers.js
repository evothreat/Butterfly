var workers;

function loadWorkers() {
    $.getJSON('/api/v1/workers', function (res) {
        workers = res.reduce(function (map, w) {            // TODO: use array with filter?
            map[w.id] = w;
            return map;
        }, {});
    });
}

function createWorkerTable() {
    if (workers == null) {
        loadWorkers();
    }
    $('#worker-table').DataTable({
            data: workers,
            columns: [
                {data: 'id', title: 'ID'},
                {data: 'hostname', title: 'Hostname'},
                {data: 'ip_addr', title: 'IP-Addr'},
                {data: 'os', title: 'OS'},
                {data: 'country', title: 'Country'},
                {data: 'last_seen', title: 'Last-Seen'},
                {
                    data: null,
                    title: 'Action',
                    render: function (data, type, row) {
                        return `<button type="button" class="btn action-btn"><i class="bi bi-box-arrow-up-left"></i></button>
                                <button type="button" class="btn action-btn" onclick="showResourceInfo(${row.id})"><i class="bi bi-info-circle-fill"></i></button>`;
                    }
                }
            ]
        }
    );
}