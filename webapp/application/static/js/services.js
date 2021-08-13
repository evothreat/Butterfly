var resourceInfo;


function loadResourceInfo() {
    $.getJSON('/api/v1/workers/-/resource-info', function (res) {
        resourceInfo = res.reduce(function (map, ri) {
            map[ri.worker_id] = ri;
            return map;
        }, {});
    });
}

function showResourceInfo(wid) {
    var ri = resourceInfo[wid];
    $('#cpu-info').html(ri.cpu);
    $('#gpu-info').html(ri.gpu);
    $('#ram-info').html(ri.ram);
    $('#resource-dlg').css('display', 'flex');
}
