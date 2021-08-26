var currWorkerId;
var jobsTable,
    uploadsTable;
// TODO: add workersTable

// BYTES TO HUMAN-READABLE
function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// MODAL
function showModal(modalId) {
    $(modalId).removeClass('hidden');
}

function hideModal(modalId) {
    $(modalId).addClass('hidden');
}

// RESOURCE INFO
function showResourceInfo(wid) {
    $.getJSON(`/api/v1/workers/${wid}/resource-info`, function (res) {
        $('#cpu-info').html(res.cpu);
        $('#gpu-info').html(res.gpu);
        $('#ram-info').html(res.ram);
        showModal('#resource-dlg');
    });
    // TODO: check for errors!
}

// TABS
function switchTab() {
    // get current tab
    var currTab = $('.tabs .active-tab');
    currTab.removeClass('active-tab');
    $('#' + currTab.attr('aria-controls')).removeClass('active-body');
    // set clicked tab
    var selected = $(this);
    selected.addClass('active-tab');
    $('#' + selected.attr('aria-controls')).addClass('active-body');
}

// TABLE ROWS
function getSelectedRows(tableId) {
    var selected = [];
    $(tableId + ' .one-select').each(function () {
        selected.push($(this).val());     // TODO: use .all-select? convert to number? avoid injection? add active class?
    });
    return selected;
}

function selectAllRows() {
    var box = $(this);
    var val = box.is(':checked');
    box.closest('table').find('.one-select').each(function () {
        $(this).prop('checked', val);
    })
}

// HISTORY
function setCurrWorkerId() {
    var url = window.location.pathname;
    currWorkerId = url.substring(url.lastIndexOf('/') + 1);
}

// JOBS
function createJobsTable() {
    jobsTable = $('#jobs-table').DataTable({
        ajax: {
            url: `/api/v1/workers/${currWorkerId}/jobs`,
            dataSrc: '',
            error: function () {
                alert('Failed to load jobs data!');
            }
        },
        columns: [
            {
                data: null,
                title: `<input class="all-select" type="checkbox">`,
                render: function (data, type, row) {
                    if (type === 'display') {
                        return `<input class="one-select" type="checkbox" value="${row.id}"/>`;
                    }
                    return null;
                }
            },
            {data: 'id', title: 'ID'},
            {data: 'todo', title: 'ToDo'},
            {data: 'completed', title: 'Completed'},
            {data: 'created', title: 'Created'},
            {
                data: null,
                title: 'Action',
                render: function (data, type, row) {
                    if (type === 'display') {
                        var res = `<button type="button" onclick="removeJob(${row.id})" class="action-btn">
                                    <i class="fa fa-trash" aria-hidden="true"></i>
                                </button>`;
                        if (row.completed) {
                            res += `<button type="button" onclick="window.open('/api/v1/workers/${currWorkerId}/jobs/${row.id}/report')" 
                                    class="action-btn"> <i class="fa fa-search" aria-hidden="true"></i>
                                </button>`;
                        }
                        return res;
                    }
                    return null;
                }
            }],
        columnDefs: [
            {
                searchable: false,
                orderable: false,
                targets: [0, 5]
            },
            {
                className: 'dt-body-center',
                targets: [0]
            }
        ],
        order: [[1, 'asc']]
    });
}

function removeJobApi(jobId, func) {
    $.ajax({
        url: `/api/v1/workers/${currWorkerId}/jobs/${jobId}`,
        type: 'DELETE',
        success: func,
        error: function () {
            alert('Failed to delete job!');
        }
    });
}

function removeJob(jobId) {
    removeJobApi(jobId, function () {
        jobsTable.rows(function (ix, data) {
            return jobId === data.id;
        }).remove().draw();
    });
}

function createJobApi(job, func) {
    $.ajax({
        type: 'POST',
        url: `/api/v1/workers/${currWorkerId}/jobs`,
        data: JSON.stringify(job),
        contentType: 'application/json',
        dataType: 'json',
        success: func,
        error: function () {
            alert('Failed to create job!');
        }
    });
}

function createJob() {
    var job = {todo: $('#todo').val()};
    createJobApi(job, function (data) {
        jobsTable.row.add(data).draw();
    })
}

// UPLOADS
function createUploadsTable() {
    uploadsTable = $('#uploads-table').DataTable({
        ajax: {
            url: `/api/v1/workers/${currWorkerId}/uploads/0/info`,
            dataSrc: '',
            error: function () {
                alert("Failed to load uploads data!");
            }
        },
        columns: [
            {
                data: null,
                title: `<input class="all-select" type="checkbox">`,
                render: function (data, type, row) {
                    if (type === 'display') {
                        return `<input class="one-select" type="checkbox" value="${row.id}"/>`
                    }
                    return null;
                }
            },
            {data: 'id', title: 'ID'},
            {data: 'filename', title: 'Filename'},
            {data: 'type', title: 'Type'},
            {
                data: 'size',
                title: 'Size',
                render: function (data, type) {
                    if (type === 'display') {
                        return formatBytes(parseInt(data));
                    }
                    return data;
                }
            },
            {data: 'created', title: 'Created'},
            {
                data: null,
                title: 'Action',
                render: function (data, type, row) {
                    if (type === 'display') {
                        return `<button type="button" onclick="removeUpload(${row.id})" class="action-btn">
                                    <i class="fa fa-trash" aria-hidden="true"></i>
                                </button>
                                <button type="button" class="action-btn"
                                    onclick="window.location.href='/api/v1/workers/${currWorkerId}/uploads/${row.id}?attach'">
                                    <i class="fa fa-download" aria-hidden="true"></i>
                                </button>
                                <button type="button" class="action-btn"
                                    onclick="window.open('/api/v1/workers/${currWorkerId}/uploads/${row.id}')">
                                    <i class="fa fa-eye" aria-hidden="true"></i>
                                </button>`;
                    }
                    return null;
                }
            }],
        columnDefs: [
            {
                searchable: false,
                orderable: false,       // exclude checkbox & buttons
                targets: [0, 6]
            },
            {
                searchable: false,
                targets: [4]            // exclude size
            },
            {
                className: 'dt-body-center',
                targets: [0]            // center checkboxes
            }
        ],
        order: [[1, 'asc']]
    });
}

function removeUploadApi(uploadId, func) {
    $.ajax({
        url: `/api/v1/workers/${currWorkerId}/uploads/${uploadId}`,
        type: 'DELETE',
        success: func,
        error: function () {
            alert('Failed to delete upload!');
        }
    });
}

function removeUpload(uploadId) {
    removeUploadApi(uploadId, function () {
        uploadsTable.rows(function (ix, data) {
            return uploadId === data.id;
        }).remove().draw();
    });
}

// TODO: instead of url use id
function retrieveReport(jobId, func, n = 10) {
    $.ajax({
        url: `/api/v1/workers/${currWorkerId}/jobs/${jobId}/report`,
        type: 'GET',
        success: func,
        error: function () {
            if (n > 0) {
                setTimeout(function () {
                    retrieveReport(jobId, func, n - 1);
                }, 5000);
            } else {
                alert('Failed to retrieve report!');
            }
        }
    })
}