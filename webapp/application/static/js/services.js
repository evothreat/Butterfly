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
                        defaultContent: '<i class="fa fa-info-circle fa-lg"></i>',
                        orderable: false
                    },
                    {
                        data: null,
                        className: 'dt-center worker-screenshot',
                        defaultContent: '<i class="fa fa-picture-o fa-lg"></i>',
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
    });
}