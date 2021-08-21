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

function createWorkersTable() {
    workersTable = $('#workers-table').DataTable({
            ajax: {
                url: '/api/v1/workers',
                dataSrc: '',
                error: function () {
                    alert("Failed to load workers data!");
                }
            },
            columns: [
                {
                    data: null,
                    title: `<input class="all-select" type="checkbox" onclick="selectAllRows('#workers-table')">`,
                    render: function (data, type, row) {
                        return '<input class="one-select" type="checkbox" value="' + row.id + '"/>'
                    }
                },
                {data: 'id', title: 'ID'},
                {data: 'hostname', title: 'Hostname'},
                {data: 'ip_addr', title: 'IP-Addr'},
                {data: 'os', title: 'OS'},
                {data: 'country', title: 'Country'},
                {
                    data: 'last_seen',
                    title: 'Last-Seen',
                    /*render: function (data, type, row) {
                        var diff = Math.abs(new Date() - new Date(data));
                    }*/
                },
                {
                    data: null,
                    title: 'Action',
                    render: function (data, type, row) {
                        return `<button type="button" class="action-btn" 
                                        onclick="document.location.href='/workers/` + row.id + `';">
                                    <i class="fa fa-pencil-square-o" aria-hidden="true"></i>
                                </button>
                                <button type="button" class="action-btn" onclick="showResourceInfo(` + row.id + `)">
                                    <i class="fa fa-info-circle" aria-hidden="true"></i>
                                </button>`;
                    }
                }
            ],
            columnDefs: [
                {
                    searchable: false,
                    orderable: false,
                    targets: [0, 7]
                },
                {
                    className: 'dt-body-center',
                    targets: [0]
                }
            ],
            order: [[1, 'asc']]
        }
    );
}