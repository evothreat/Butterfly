let currWorkerId;
let workersTable,
    jobsTable,
    uploadsTable;

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
    $.getJSON(`/api/workers/${wid}/resource-info`, function (res) {
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
    let currTab = $('.tabs .active-tab');
    currTab.removeClass('active-tab');
    $('#' + currTab.attr('aria-controls')).removeClass('active-body');
    // set clicked tab
    let selected = $(this);
    selected.addClass('active-tab');
    $('#' + selected.attr('aria-controls')).addClass('active-body');
}

// TABLE ROWS
function getSelectedRows(tableId) {
    let selected = [];
    $(tableId + ' .one-select').each(function () {
        selected.push($(this).val());     // TODO: use .all-select? convert to number? avoid injection? add active class?
    });
    return selected;
}

function selectAllRows() {
    let box = $(this);
    let val = box.is(':checked');
    box.closest('table').find('.one-select').each(function () {
        $(this).prop('checked', val);
    })
}

// WORKERS
function createWorkersTable() {
    workersTable = $('#workers-table').DataTable({
            ajax: {
                url: '/api/workers',
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
                        if (type === 'display') {
                            return `<input class="one-select" type="checkbox" value="${row.id}"/>`;
                        }
                        return null;
                    }
                },
                {data: 'hostname', title: 'Hostname'},
                {data: 'ip_addr', title: 'IP-Address'},
                {data: 'country', title: 'Country'},
                {data: 'os', title: 'OS'},
                {
                    data: 'is_admin',
                    title: 'Admin',
                    render: function (data, type) {
                        if (type === 'display' || type === 'filter') {
                            return data ? 'yes' : 'no';
                        }
                        return data;
                    }
                },
                {
                    data: "boost",
                    title: "Boost",
                    render: function (data, type) {
                        if (type === 'display' || type === 'filter') {
                            return data ? 'on' : 'off';
                        }
                        return data;
                    }
                },
                {
                    data: 'last_seen',
                    title: 'Last-Seen',
                    render: function (data, type) {
                        if (type === 'display') {
                            let diff = new Date() - new Date(data);
                            if (diff > 60000) {
                                return data.slice(0, 19).replace('T', ' ');
                            }
                            return 'Online';
                        }
                        return data;
                    }
                },
                {
                    data: null,
                    title: 'Action',
                    render: function (data, type, row) {
                        if (type === 'display') {
                            return `<button type="button" class="action-btn" 
                                            onclick="document.location.href='/workers/${row.id}?boost=${row.boost}'">
                                        <i class="fa fa-pencil-square-o" aria-hidden="true"></i>
                                    </button>
                                    <button type="button" class="action-btn" onclick="showResourceInfo('${row.id}')">
                                        <i class="fa fa-info-circle" aria-hidden="true"></i>
                                    </button>`;
                        }
                        return null;
                    }
                }
            ],
            columnDefs: [
                {
                    searchable: false,
                    orderable: false,
                    targets: [0, 8]
                },
                {
                    className: 'dt-body-center',
                    targets: [0]
                }
            ],
            order: [[7, 'asc']]
        }
    );
}

// HISTORY
function setCurrWorkerId() {
    let url = window.location.pathname;
    currWorkerId = url.substring(url.lastIndexOf('/') + 1);
}

// JOBS
function createJobsTable() {
    jobsTable = $('#jobs-table').DataTable({
        ajax: {
            url: `/api/workers/${currWorkerId}/jobs`,
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
            {
                data: 'is_done',
                title: 'Completed',
                render: function (data, type) {
                        if (type === 'display' || type === 'filter') {
                            return data ? 'yes' : 'no';
                        }
                        return data;
                    }
            },
            {
                data: 'created',
                title: 'Created',
                render: function (data, type) {
                        if (type === 'display') {
                                return data.slice(0, 19).replace('T', ' ');
                        }
                        return data;
                    }
            },
            {
                data: null,
                title: 'Action',
                render: function (data, type, row) {
                    if (type === 'display') {
                        let res = `<button type="button" onclick="removeJob(${row.id})" class="action-btn">
                                        <i class="fa fa-trash" aria-hidden="true"></i>
                                   </button>`;
                        if (row.is_done) {
                            res += `<button type="button" onclick="window.open('/api/workers/${currWorkerId}/jobs/${row.id}/report')" 
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
        url: `/api/workers/${currWorkerId}/jobs/${jobId}`,
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
        url: `/api/workers/${currWorkerId}/jobs`,
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
    let job = {todo: $('#todo').val()};
    createJobApi(job, function (data) {
        jobsTable.row.add(data).draw();
    })
}

// UPLOADS
function createUploadsTable() {
    uploadsTable = $('#uploads-table').DataTable({
        ajax: {
            url: `/api/workers/${currWorkerId}/uploads/0/info`,
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
            {
                data: 'created',
                title: 'Created',
                render: function (data, type) {
                        if (type === 'display') {
                                return data.slice(0, 19).replace('T', ' ');
                        }
                        return data;
                    }
            },
            {
                data: null,
                title: 'Action',
                render: function (data, type, row) {
                    if (type === 'display') {
                        return `<button type="button" onclick="removeUpload(${row.id})" class="action-btn">
                                    <i class="fa fa-trash" aria-hidden="true"></i>
                                </button>
                                <button type="button" class="action-btn"
                                        onclick="window.location.href='/api/workers/${currWorkerId}/uploads/${row.id}?attach'">
                                    <i class="fa fa-download" aria-hidden="true"></i>
                                </button>
                                <button type="button" class="action-btn"
                                        onclick="window.open('/api/workers/${currWorkerId}/uploads/${row.id}')">
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
        url: `/api/workers/${currWorkerId}/uploads/${uploadId}`,
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

function retrieveReport(jobId, funSucc, funcErr, n = 10) {
    $.ajax({
        url: `/api/workers/${currWorkerId}/jobs/${jobId}/report`,
        type: 'GET',
        success: funSucc,
        error: function (xhr, stat, error) {
            if (n > 1) {
                setTimeout(function () {
                    retrieveReport(jobId, funSucc, funcErr, n - 1);
                }, 5000);
            } else {
                funcErr(xhr, stat, error);
            }
        }
    })
}

function setupBoostToggle() {
    let urlSearch = new URLSearchParams(window.location.search);
    let boost = $('#boost');

    boost.prop('checked', urlSearch.get('boost') === 'true');
    boost.change(function () {
        createJobApi({                                          // TODO: call directly createJob() and pass job?
            todo: 'boost ' + ($(this).is(':checked') ? 'on' : 'off')
        }, null);
    })
}