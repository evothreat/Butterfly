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
    $(modalId).css('display', 'flex');
}

function hideModal(modalId) {
    $(modalId).css('display', 'none')
}

// RESOURCE INFO
function showResourceInfo(wid) {
    $.getJSON('/api/v1/workers/' + wid + '/resource-info', function (res) {
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
        selected.push($(this).val());               // TODO: use .all-select? convert to number? avoid injection? add active class?
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
            url: '/api/v1/workers/' + currWorkerId + '/jobs',
            dataSrc: '',
            error: function () {
                alert("Failed to load jobs data!");
            }
        },
        columns: [
            {
                data: null,
                title: `<input class="all-select" type="checkbox">`,
                render: function (data, type, row) {
                    return `<input class="one-select" type="checkbox" value="${row.id}"/>`
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
                    return `<button type="button" onclick="removeJob(${row.id})" class="action-btn">
                                    <i class="fa fa-trash" aria-hidden="true"></i>
                            </button>`;
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

function removeJob(jobId) {
    $.ajax({
        url: '/api/v1/workers/' + currWorkerId + '/jobs/' + jobId,
        type: 'DELETE',
        success: function () {
            jobsTable.rows(function (ix, data) {
                return jobId === data.id;
            }).remove().draw();
        },
        error: function () {
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
        success: function (data) {
            jobsTable.row.add(data).draw();
        },
        error: function () {
            alert('Failed to create job!');
        }
    });
}

// UPLOADS
function createUploadsTable() {
    uploadsTable = $('#uploads-table').DataTable({
        ajax: {
            url: '/api/v1/workers/' + currWorkerId + '/uploads/0/info',
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
                    return `<input class="one-select" type="checkbox" value="${row.id}"/>`
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

function removeUpload(uploadId) {
    $.ajax({
        url: '/api/v1/workers/' + currWorkerId + '/uploads/' + uploadId,
        type: 'DELETE',
        success: function () {
            uploadsTable.rows(function (ix, data) {
                return uploadId === data.id;
            }).remove().draw();
        },
        error: function () {
            alert('Failed to delete upload!');
        }
    });
}