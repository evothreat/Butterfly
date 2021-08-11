function createWorkerTable() {
    $.getJSON('/api/v1/workers', function (res) {     // TODO: handle errors
        var workerTbl = $('#worker-table');
        var workerDataTbl = workerTbl.DataTable({
                data: res,
                columns: [
                    {data: 'id', title: 'ID'},
                    {data: 'hostname', title: 'Hostname'},
                    {data: 'ip_addr', title: 'IP-Addr'},
                    {data: 'os', title: 'OS'},
                    {data: 'country', title: 'Country'},
                    {data: 'last_seen', title: 'Last-Seen'},
                    {
                        data: null,
                        className: 'dt-center worker-info',
                        defaultContent: '<span style="cursor: pointer" class="fa fa-info-circle fa-lg"></span>',
                        orderable: false
                    },
                    {
                        data: null,
                        className: 'dt-center worker-screenshot',
                        defaultContent: '<span style="cursor: pointer" class="fa fa-picture-o fa-lg"></span>',
                        orderable: false
                    }
                ]
            }
        );
        workerTbl.on('click', 'td.worker-info', function (e) {
            //e.preventDefault();
            //console.log(workerDataTbl.row(this).data());
            // TODO: redirect to new page
        });
        workerTbl.on('click', 'td.worker-screenshot', function (e) {
            //e.preventDefault();
            //console.log(workerDataTbl.row(this).data());
            // TODO: redirect to new page
        });
    });
}

function createWorkerTable2() {
    $.getJSON('/api/v1/workers', function (res) {     // TODO: handle errors
         $('#worker-table').DataTable({
                data: res,
                columns: [
                    {data: 'id', title: 'ID'},
                    {data: 'hostname', title: 'Hostname'},
                    {data: 'ip_addr', title: 'IP-Addr'},
                    {data: 'os', title: 'OS'},
                    {data: 'country', title: 'Country'},
                    {data: 'last_seen', title: 'Last-Seen'},
                    {
                        data: null,
                        title: 'Operation',
                        render: function (data, type, row) {
                            return `<a style="color: #343a40" href="/workers/${row.id}"><span class="bi bi-info-circle-fill worker-info"></span></a>
                                    <a style="color: #343a40" href="/workers/${row.id}"><span class="bi bi-camera-fill worker-screenshot"></span></a>`
                        }
                    }
                ]
            }
        );
    });
}