from sqlalchemy import and_
from sqlalchemy.exc import IntegrityError
from flask import jsonify, request, send_file
from re import match
from os import mkdir, remove as remove_file
from os.path import join as path_join, getsize, getctime, splitext
from application import *
from application.models import *


# TODO: handle wrong type errors (in json)!


def valid_filename(name):
    return match('^[\w\-. ]+$', name)


# WORKERS -----------------------------------------------------------------------

@app.route('/api/v1/workers', methods=['GET'])
def get_workers():
    return jsonify(Worker.query.all()), 200


@app.route('/api/v1/workers', methods=['POST'])
def create_worker():
    if not request.is_json:
        return '', 400
    try:
        w = Worker.from_dict(request.json)
    except KeyError:
        return '', 422
    db.session.add(w)
    db.session.commit()
    # Create upload dir
    mkdir(path_join(app.config['UPLOADS_DIR'], w.id))
    return jsonify(w), 201


@app.route('/api/v1/workers/<wid>', methods=['GET'])
def get_single_worker(wid):
    w = Worker.query.get(wid)
    return (jsonify(w), 200) if w else ('', 404)                                        # TODO: use get_or_404 instead?


@app.route('/api/v1/workers/<wid>', methods=['DELETE'])
def delete_worker(wid):
    if Worker.query.filter_by(id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


# JOBS -----------------------------------------------------------------------

@app.route('/api/v1/workers/<wid>/jobs', methods=['GET'])
def get_jobs(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    if 'done' in request.args:
        jobs = Job.query.filter_by(worker_id=wid, completed=True).all()
    elif 'undone' in request.args:
        jobs = Job.query.filter_by(worker_id=wid, completed=False).all()
    else:
        jobs = Job.query.filter_by(worker_id=wid).all()
    return jsonify(jobs), 200


@app.route('/api/v1/workers/<wid>/jobs', methods=['POST'])
def create_job(wid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):                                    # TODO: remove after enabling foreign keys?
        return '', 404
    try:
        job = Job.from_dict(request.json)
        job.worker_id = wid
    except KeyError:
        return '', 422
    db.session.add(job)
    db.session.commit()
    return jsonify(job), 201


@app.route('/api/v1/workers/<wid>/jobs/<int:jid>', methods=['DELETE'])
def delete_job(wid, jid):
    if Job.query.filter_by(id=jid, worker_id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<wid>/jobs/<int:jid>', methods=['PATCH'])
def update_job(wid, jid):
    if not request.is_json:
        return '', 400
    job = Job.query.filter_by(id=jid, worker_id=wid).first()
    if not job:
        return '', 404
    # job.id = request.json.get('id', job.id)
    # job.todo = request.json.get('todo', job.todo)
    job.completed = request.json.get('completed', job.completed)
    # job.created = request.json.get('created', job.created)
    # job.report_type = request.json.get('report_type', job.report_type)
    # job.worker_id = request.json.get('worker_id', job.worker_id)
    db.session.commit()
    return '', 200


# RESOURCE INFO --------------------------------------------------------------------
@app.route('/api/v1/workers/<wid>/resource-info', methods=['POST'])
def create_resource_info(wid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):
        return '', 404
    try:
        ri = ResourceInfo.from_dict(request.json)
        ri.worker_id = wid
        db.session.add(ri)
        db.session.commit()
    except KeyError:
        return '', 422
    except IntegrityError:
        db.session.rollback()
        return '', 409
    return jsonify(ri), 201


@app.route('/api/v1/workers/<wid>/resource-info', methods=['GET'])          # TODO: add new path /workers/resource-info
def get_resource_info(wid):
    if wid == '-':
        return jsonify(ResourceInfo.query.all()), 200
    ri = ResourceInfo.query.get(wid)
    return (jsonify(ri), 200) if ri else ('', 404)


# UPLOADS --------------------------------------------------------------------
@app.route('/api/v1/workers/<wid>/uploads', methods=['POST'])
def create_upload(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    file = request.files.get('file')  # name attribute in upload html
    if not file or not valid_filename(file.filename):
        return '', 422
    filename = file.filename
    _, ext = splitext(filename)
    ext = ext[1:].upper() if ext else 'NONE'
    path = path_join(app.config['UPLOADS_DIR'], wid, filename)
    file.save(path)                                                                        # TODO: handle failed to save
    file.close()
    up = Upload(filename=filename, size=getsize(path),
                type=ext, created=datetime.fromtimestamp(getctime(path)),
                worker_id=wid)
    db.session.add(up)
    db.session.commit()
    return jsonify(up), 201


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>', methods=['GET'])
def get_single_upload(wid, uid):
    up = Upload.query.filter_by(id=uid, worker_id=wid).first()
    if not up:
        return '', 404
    filepath = path_join(app.config['UPLOADS_DIR'], wid, up.filename)
    attach = bool('attach' in request.args)
    try:
        return send_file(filepath, as_attachment=attach)
    except FileNotFoundError:
        return '', 404  # remove from database too


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>', methods=['DELETE'])           # TODO: add own parameter decoder?
def delete_upload(wid, uid):
    up = Upload.query.filter_by(id=uid, worker_id=wid).first()
    if not up:
        return '', 404
    filepath = path_join(app.config['UPLOADS_DIR'], wid, up.filename)
    try:
        remove_file(filepath)
    except FileNotFoundError:
        pass  # or return '', 404
    db.session.delete(up)
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>/info', methods=['GET'])
def get_upload_info(wid, uid):
    if uid == 0:
        return jsonify(Upload.query.filter_by(worker_id=wid).all()), 200
    up = Upload.query.filter_by(worker_id=wid, id=uid)
    return (jsonify(up), 200) if up else ('', 404)


# REPORTS --------------------------------------------------------------------
@app.route('/api/v1/workers/<wid>/jobs/<int:jid>/report', methods=['POST'])
def create_report(wid, jid):
    if not obj_exists(and_(Job.worker_id == wid, Job.id == jid)):
        return '', 404
    if request.mimetype != 'text/plain' or request.content_length > app.config['REPORT_LEN']:
        return '', 422
    rep = JobReport(job_id=jid, report=request.get_data(cache=False))
    try:
        db.session.add(rep)
        Job.query.filter_by(id=jid).update({'completed': True})
        db.session.commit()
    except IntegrityError:  # report already exists, means Job.completed == True
        db.session.rollback()
        return '', 409
    return '', 201


@app.route('/api/v1/workers/<wid>/jobs/<int:jid>/report', methods=['GET'])
def get_report(wid, jid):
    if not obj_exists(and_(Job.worker_id == wid, Job.id == jid)):
        return '', 404
    rep = JobReport.query.get(jid)
    return (rep.report, 200) if rep else ('', 404)
