function createWorkerTable() {
    console.log("HELLO");
    $.getJSON('/api/v1/workers', function (res) {                                                                    // TODO: handle errors
        res.forEach(w => w['operation'] = '<a href="#" class="op-link"><i class="fa fa-info-circle fa-lg"></i></a>');
        for (var i = 0; i < res.length; i++) {
            id_link = '#' + res[i]['id'];                                                                                  // TODO: add url!
            res[i]['operation'] = '<a href="' + id_link + '" class="op-link"><i class="fa fa-info-circle fa-lg"></i></a> \
                                   <a href="' + id_link + '" class="op-link"><i class="fa fa-picture-o fa-lg"></i></a>';
        }
        console.log(res);
        $('#worker-table').DataTable({
                data: res,
                columns: [
                    {'data': 'id', 'title': 'ID'},
                    {'data': 'hostname', 'title': 'Hostname'},
                    {'data': 'ip_addr', 'title': 'IP-Addr'},
                    {'data': 'os', 'title': 'OS'},
                    {'data': 'country', 'title': 'Country'},
                    {'data': 'last_seen', 'title': 'Last-Seen'},
                    {'data': 'operation', 'title': 'Operation'}]
            }
        )
    });
}