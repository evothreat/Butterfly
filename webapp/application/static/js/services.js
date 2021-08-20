var resourceInfo;
var currWorkerId;
var jobsTable;

// MODAL
function showModal(modalId) {
    $(modalId).css('display', 'flex');
}

function hideModal(modalId) {
    $(modalId).css('display', 'none')
}

// RESOURCE INFO
function loadResourceInfo() {
    $.getJSON('/api/v1/workers/0/resource-info', function (res) {
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
    showModal('#resource-dlg');
}

// TABS
function setActiveTab(tabId) {
    $('.tabs .active-tab').removeClass('active-tab');
    $(tabId).addClass('active-tab');
}

function loadContent(tabId) {
    /*switch (tabId) {
        case '#history-tab':
        case '#downloads-tab':
        case '#manager-tab':
        case '#terminal-tab':
        case '#credentials-tab':
        case '#cronjob-tab':
        case '#sniffer-tab':
    }*/
    setActiveTab(tabId);
}

// TABLE ROWS
function getSelectedRows(tableId) {
    var selected = [];
    $(tableId + ' .one-select').each(function () {
        selected.push($(this).val());               // TODO: use .all-select? convert to number? avoid injection?
    });
    return selected;
}

function selectAllRows(tableId) {
    var current = $(tableId + ' .all-select').is(':checked');
    $(tableId + ' .one-select').each(function () {
        $(this).prop('checked', current);
    })
}

// HISTORY
function setCurrWorkerId() {
    var url = window.location.pathname;
    currWorkerId = url.substring(url.lastIndexOf('/') + 1);
}

function removeJob(jobId) {
    $.ajax({
        url: '/api/v1/workers/' + currWorkerId + '/jobs/' + jobId,
        type: 'DELETE',
        success: function () {
            jobsTable.rows(function (ix, data, node) {
                return jobId === data.id;
            }).remove().draw();
        },
        error: function (xhr, stat, err) {
            alert('Failed to delete job!');
        }
    });
}

function createJob() {
    var job = {}
    job.todo = $('#todo').val();
    $.ajax({
        type: "POST",
        url: '/api/v1/workers/' + currWorkerId + '/jobs',
        data: JSON.stringify(job),
        contentType: 'application/json',
        dataType: 'json',
        success: function (data, stat, xhr) {
            data.checkbox = '<td class="dt-body-center"><input class="one-select" type="checkbox" value="{{ j.id }}"/></td>'
            data.action = '<button type="button" onclick="removeJob(' + data.id + ')" class="action-btn"> \
                           <i class="fa fa-trash" aria-hidden="true"></i></button>';
            jobsTable.row.add(data).draw();
        },
        error: function (xhr, stat, err) {
            alert('Failed to create job!');
        }
    });

}