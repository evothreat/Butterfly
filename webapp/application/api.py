from sqlalchemy.exc import IntegrityError
from flask import jsonify, request, send_file
from re import match
from os import mkdir, remove as remove_file
from os.path import join as path_join, getsize, getctime, splitext
from mimetypes import guess_extension
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
        db.session.add(w)
        db.session.commit()
        # Creating uploads directory
        mkdir(path_join(app.config['UPLOADS_DIR'], str(w.id)))
        mkdir(path_join(app.config['REPORTS_DIR'], str(w.id)))
    except KeyError:
        return '', 422
    except IntegrityError:
        return '', 409
    return jsonify(w), 201


@app.route('/api/v1/workers/<wid>', methods=['GET'])
def get_single_worker(wid):
    w = Worker.query.get(wid)
    return ('', 404) if not w else (jsonify(w), 200)  # TODO: use get_or_404 instead


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
    arg = request.args.get('is_done', '').lower()
    if arg and (arg == 'true' or arg == 'false'):
        jobs = Job.query.filter_by(worker_id=wid, is_done=bool(arg == 'true')).all()
    else:
        jobs = Job.query.filter_by(worker_id=wid).all()
    return jsonify(jobs), 200


@app.route('/api/v1/workers/<wid>/jobs', methods=['POST'])
def create_job(wid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):  # TODO: remove after enabling foreign keys?
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
    if not obj_exists(Worker.id == wid):
        return '', 404
    if Job.query.filter_by(id=jid, worker_id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<wid>/jobs/<int:jid>', methods=['PATCH'])
def update_job(wid, jid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):
        return '', 404
    job = Job.query.filter_by(id=jid, worker_id=wid).first()
    if not job:
        return '', 404
    # job.id = request.json.get('id', job.id)
    # job.todo = request.json.get('todo', job.todo)
    job.is_done = request.json.get('is_done', job.is_done)
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
    except KeyError:
        return '', 422
    db.session.add(ri)
    db.session.commit()
    return jsonify(ri), 201


@app.route('/api/v1/workers/<wid>/resource-info', methods=['GET'])  # TODO: add new path /workers/resource-info
def get_resource_info(wid):
    if wid == '-':
        return jsonify(ResourceInfo.query.all()), 200
    if obj_exists(Worker.id == wid):
        return jsonify(ResourceInfo.query.filter_by(worker_id=wid).first()), 200
    return '', 404


# UPLOADS --------------------------------------------------------------------
@app.route('/api/v1/workers/<wid>/uploads', methods=['POST'])
def create_upload(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    file = request.files.get('file')                    # name attribute in upload html
    if not file or not valid_filename(file.filename):
        return '', 422
    filename = file.filename
    _, ext = splitext(filename)
    ext = ext[1:].upper() if ext else 'NONE'
    path = path_join(app.config['UPLOADS_DIR'], wid, filename)
    file.save(path)
    file.close()
    up = Upload(filename=filename, size=getsize(path),
                type=ext, created=datetime.fromtimestamp(getctime(path)),
                worker_id=wid)
    db.session.add(up)
    db.session.commit()
    return jsonify(up), 201


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>', methods=['GET'])
def get_single_upload(wid, uid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    up = Upload.query.filter_by(id=uid, worker_id=wid).first()
    if not up:
        return '', 404
    filepath = path_join(app.config['UPLOADS_DIR'], wid, up.filename)
    attach = bool('attach' in request.args)
    try:
        return send_file(filepath, as_attachment=attach)
    except FileNotFoundError:
        return '', 404


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>', methods=['DELETE'])  # TODO: add own parameter decoder?
def delete_upload(wid, uid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    up = Upload.query.filter_by(id=uid, worker_id=wid).first()
    if not up:
        return '', 404
    filepath = path_join(app.config['UPLOADS_DIR'], wid, up.filename)
    try:
        remove_file(filepath)
    except FileNotFoundError:
        return '', 404
    db.session.delete(up)
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<wid>/uploads/<int:uid>/info', methods=['GET'])
def get_upload_info(wid, uid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    if uid == 0:
        return jsonify(Upload.query.filter_by(worker_id=wid).all()), 200
    return jsonify(Upload.query.filter_by(id=uid, worker_id=wid).first()), 200


# REPORTS --------------------------------------------------------------------
@app.route('/api/v1/workers/<wid>/jobs/<int:jid>/report', methods=['POST'])
def create_report(wid, jid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    job = Job.query.filter_by(id=jid, worker_id=wid).first()
    if not job:
        return '', 404
    if not job.is_done or job.report_type:          # if job undone or already reported
        return '', 409
    path = path_join(app.config['REPORTS_DIR'], wid, 'report_{}{}')
    if request.mimetype == 'multipart/form-data':
        file = request.files.get('file')
        if not file:
            return '', 422
        _, ext = splitext(file.filename)            # or use guess_extension(file.content_type)
        file.save(path.format(jid, ext))            # TODO: if no extension, return error
        file.close()
    else:
        ext = guess_extension(request.content_type)
        with open(path.format(jid, ext), 'wb') as f:
            f.write(request.get_data(cache=False))
    job.report_type = ext
    db.session.commit()
    return '', 204
