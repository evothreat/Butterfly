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
        error: function (xhr, stat, err) {
            alert('Failed to delete job!');
        }
    });
    jobsTable.rows(function (ix, data, node) {
        return jobId === data[1];
    }).remove().draw();
}