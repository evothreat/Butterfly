from sqlalchemy.exc import IntegrityError
from flask import jsonify, request, send_file
from werkzeug.utils import secure_filename
from os import mkdir
from os.path import join as path_join, getsize, getctime, splitext
import datetime
from application import *
from application.models import *


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
    except KeyError:
        return '', 422
    except IntegrityError:
        return '', 409
    return jsonify(w), 201


@app.route('/api/v1/workers/<int:wid>', methods=['GET'])
def get_single_worker(wid):
    w = Worker.query.get(wid)
    return ('', 404) if not w else (jsonify(w), 200)  # TODO: use get_or_404 instead


@app.route('/api/v1/workers/<int:wid>', methods=['DELETE'])
def delete_worker(wid):
    if Worker.query.filter_by(id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


# JOBS -----------------------------------------------------------------------

@app.route('/api/v1/workers/<int:wid>/jobs', methods=['GET'])
def get_jobs(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    args = request.args.get('is_done')
    is_done = bool(args and args.lower() == 'true')
    jobs = Job.query.filter_by(worker_id=wid, is_done=is_done).all()
    return jsonify(jobs), 200


@app.route('/api/v1/workers/<int:wid>/jobs', methods=['POST'])
def create_job(wid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):  # TODO: remove after enabling foreign keys?
        return '', 404
    try:
        j = Job.from_dict(request.json)
        j.worker_id = wid
    except KeyError:
        return '', 422
    db.session.add(j)
    db.session.commit()
    return jsonify(j), 201


@app.route('/api/v1/workers/<int:wid>/jobs/<int:jid>', methods=['DELETE'])
def delete_job(wid, jid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    if Job.query.filter_by(id=jid, worker_id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


# OTHER --------------------------------------------------------------------
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


@app.route('/api/v1/workers/<int:wid>/resource-info', methods=['GET'])      # TODO: add new path /workers/resource-info
def get_resource_info(wid):
    if wid == 0:
        return jsonify(ResourceInfo.query.all()), 200
    if obj_exists(Worker.id == wid):
        return jsonify(ResourceInfo.query.filter_by(worker_id=wid).all()), 200
    return '', 404


@app.route('/api/v1/workers/<int:wid>/uploads', methods=['POST'])
def create_upload(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    try:  # TODO: use if-else?
        file = request.files['file']  # TODO: name attribute in upload html
    except KeyError:
        return '', 422
    name = secure_filename(file.filename)
    path = path_join(app.config['UPLOADS_DIR'], str(wid), name)
    _, ext = splitext(name)
    ext = ext[1:].upper() if ext else 'NONE'
    file.save(path)
    file.close()
    up = Upload(filename=name, size=getsize(path),
                type=ext, created=datetime.datetime.fromtimestamp(getctime(path)),
                worker_id=wid)
    db.session.add(up)
    db.session.commit()
    return jsonify(up), 201


@app.route('/api/v1/workers/<int:wid>/uploads/<int:uid>', methods=['GET'])
def get_single_upload(wid, uid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    up = Upload.query.get(uid)
    if not up:
        return '', 404
    filepath = path_join(app.config['UPLOADS_DIR'], str(wid), up.filename)
    try:
        return send_file(filepath)
    except FileNotFoundError:
        return '', 404


@app.route('/api/v1/workers/<int:wid>/uploads/<int:uid>/info', methods=['GET'])
def get_upload_info(wid, uid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    if uid == 0:
        return jsonify(Upload.query.filter_by(worker_id=wid).all()), 200
    return jsonify(Upload.query.filter_by(id=uid, worker_id=wid).all()), 200
