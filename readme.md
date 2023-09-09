# Butterfly

Software, that allows you  to control remote PCs using HTTP/S protocol.

### Project structure

-   WebAppGo - Contains REST API and CNC (Command & Control) written in golang.
-   WebAppPy - Contains REST API and CNC  written in python.
-   Worker - Contains client that will be started on remote PC. Written in golang.

### How it works?

#### From cnc side

1.  Client sends REST-API request to server.
2.  Server identifies request and performs operation on database.
3.  JavaScript updates the UI if necessary.

#### From client (worker) side

1.  Every n-seconds client makes REST-API request to the server. Server responses with list of jobs.
2.  Client sorts the list of jobs on creation time (ascending).
3.  Client starts to parse every job and execute it.
4. If execution fails, the error is reported else execution output/success message.

### REST-API architecture

Root path is /api

#### Workers

| Method | Path          | Description                                                          |
| ------ | ------------- | -------------------------------------------------------------------- |
| GET    | /workers      | List all workers.                                                    |
| POST   | /workers      | Create new worker.                                                   |
| GET    | /workers/:wid | Get worker or worker fields (f.e. fields=id,hostname) with given id. |
| DELETE | /workers/:wid | Delete worker with given id.                                         |
| PATCH  | /workers/:wid | Change field(s) of worker with given id.                             |

#### Jobs

| Method | Path                    | Description                                                |
| ------ | ----------------------- | ---------------------------------------------------------- |
| GET    | /workers/:wid/jobs      | List all jobs or only done/undone jobs (?done or ?undone). |
| POST   | /workers/:wid/jobs      | Create new job.                                            |
| GET    | /workers/:wid/jobs/:jid | Get job with given id.                                     |
| DELETE | /workers/:wid/jobs/:jid | Delete job with given id.                                  |

#### Hardware

| Method | Path                   | Description      |
| ------ | ---------------------- | ---------------- |
| POST   | /workers/:wid/hardware | Create hardware. |
| GET    | /workers/:wid/hardware | Delete hardware. |

#### Uploads

| Method | Path                                                                        | Description                  |
| ------ | --------------------------------------------------------------------------- | ---------------------------- |
| POST   | /workers/:wid/uploads                                                       | Create upload.               |
| GET    | /workers/:wid/uploads/:uid                                                  | Get upload with given id.    |
| DELETE | /workers/:wid/uploads/:uid                                                  | Delete upload with given id. |
| GET    | /workers/:wid/uploads/info                                                  | Get upload informations.     |
| GET    | /workers/:wid/uploads/:uid/info | Get information of upload with given id.       |

#### Reports

| Method | Path                           | Description                  |
| ------ | ------------------------------ | ---------------------------- |
| POST   | /workers/:wid/jobs/:jid/report | Create new report.           |
| GET    | /workers/:wid/jobs/:jid/report | Get report with given id.    |
| DELETE | /workers/:wid/jobs/:jid/report | Delete report with given id. |

### Features

-   Safe authentication.
-   Execute shell command and retrieve output.
-   Job execution history.
-   Upload files to remote host.
-   Download files from remote host.
-   Boost worker, so following requests are made faster.
-   Take screenshot and upload to server.
-   Show hardware information of remote host.
-   Download program and run it on remote host.
-   Show message box on remote host.

### ToDo
-   File explorer.
-   Program explorer.
-   Multi-user control (via. "owner"-field in workers table).
-   Show Network Packets in real time.

### Screenshots

#### Workers

![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/workers.png)

#### Hardware Info
![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/hardware_info.png)
#### History
![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/history.png)
#### Filter jobs
![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/filter.png)
#### Uploads
![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/uploads.png)
#### Terminal
![alt text](https://gitlab.com/evothreat/butterfly/-/raw/master/screenshots/terminal.png)
