from sqlalchemy import and_
from sqlalchemy.exc import IntegrityError
from flask import jsonify, request, send_file, url_for
from re import match
from os import mkdir, remove as remove_file
from os.path import join as path_join, exists as file_exists, getsize, getctime, splitext, dirname
from time import time

from application import *
from application.models import *


# TODO: handle wrong type errors (in json)!


def valid_filename(name):
    return match('^[\w\-. ]+$', name)


# WORKERS -----------------------------------------------------------------------

@app.route('/api/workers', methods=['GET'])
def get_workers():
    return jsonify(Worker.query.all()), 200


@app.route('/api/workers', methods=['POST'])
def create_worker():
    if not request.is_json:
        return '', 400
    try:
        w = Worker.from_dict(request.json)
        db.session.add(w)
        db.session.commit()
    except KeyError:
        return '', 422
    except IntegrityError:
        return '', 409
    mkdir(path_join(app.config['UPLOADS_DIR'], w.id))
    return '', 201, {'Location': url_for('get_single_worker', wid=w.id)}


@app.route('/api/workers/<wid>', methods=['GET'])   # TODO: add specific fields
def get_single_worker(wid):
    w = Worker.query.get(wid)
    if not w:
        return '', 404
    props = request.args.get('props')
    if props:
        try:
            data = {v: getattr(w, v) for v in props.split(',')}        # query string can contain whitespaces
        except AttributeError:
            return '', 422
        return jsonify(data), 200
    return jsonify(w), 200


@app.route('/api/workers/<wid>', methods=['DELETE'])
def delete_worker(wid):
    if Worker.query.filter_by(id=wid).delete(synchronize_session=False) == 0:
        return '', 404
    db.session.commit()
    remove_file(path_join(app.config['UPLOADS_DIR'], wid))
    return '', 200


@app.route('/api/workers/<wid>', methods=['PATCH'])
def update_worker(wid):
    if not request.is_json or not request.json:
        return '', 400
    # if Worker.query.filter_by(id=wid).update(request.json, synchronize_session=False) == 0:
    #    return '', 404
    w = Worker.query.get(wid)
    if not w:
        return '', 404
    # # w.ip_addr = request.json.get('ip_addr', w.ip_addr)
    # # w.country = request.json.get('country', w.country)
    # # w.is_admin = request.json.get('is_admin', w.is_admin)
    w.boost = request.json.get('boost', w.boost)
    db.session.commit()
    return '', 200


# JOBS -----------------------------------------------------------------------
@app.route('/api/workers/<wid>/jobs', methods=['GET'])
def get_jobs(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    return jsonify(Job.query.filter_by(worker_id=wid).all()), 200


@app.route('/api/workers/<wid>/jobs/undone', methods=['GET'])
def get_undone_jobs(wid):
    if Worker.query.filter_by(id=wid). \
            update({'last_seen': datetime.now()}, synchronize_session=False) == 0:
        return '', 404
    db.session.commit()
    return jsonify(Job.query.filter_by(worker_id=wid, is_done=False).all()), 200


@app.route('/api/workers/<wid>/jobs', methods=['POST'])
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


@app.route('/api/workers/<wid>/jobs/<int:jid>', methods=['GET'])
def get_single_job(wid, jid):
    job = Job.query.filter_by(id=jid, worker_id=wid).first()
    return (jsonify(job), 200) if job else ('', 404)


@app.route('/api/workers/<wid>/jobs/<int:jid>', methods=['DELETE'])
def delete_job(wid, jid):
    if Job.query.filter_by(id=jid, worker_id=wid).delete(synchronize_session=False) == 0:
        return '', 404
    db.session.commit()
    return '', 200


# RESOURCE INFO --------------------------------------------------------------------
@app.route('/api/workers/<wid>/hardware', methods=['POST'])
def create_hardware_info(wid):
    if not request.is_json:
        return '', 400
    if not obj_exists(Worker.id == wid):
        return '', 404
    try:
        ri = HardwareInfo.from_dict(request.json)
        ri.worker_id = wid
        db.session.add(ri)
        db.session.commit()
    except KeyError:
        return '', 422
    except IntegrityError:
        db.session.rollback()
        return '', 409
    return '', 201, {'Location': url_for('get_hardware_info', wid=wid)}


@app.route('/api/workers/<wid>/hardware', methods=['GET'])  # TODO: add new path /workers/hardware
def get_hardware_info(wid):
    if wid == '-':
        return jsonify(HardwareInfo.query.all()), 200
    ri = HardwareInfo.query.get(wid)
    return (jsonify(ri), 200) if ri else ('', 404)


# UPLOADS --------------------------------------------------------------------
@app.route('/api/workers/<wid>/uploads', methods=['POST'])
def create_upload(wid):
    if not obj_exists(Worker.id == wid):
        return '', 404
    file = request.files.get('file')                                            # name attribute in upload html
    if not file or not valid_filename(file.filename):
        return '', 422
    filename = file.filename
    name, ext = splitext(filename)
    filetype = ext[1:].upper() if ext else 'NONE'
    path = path_join(app.config['UPLOADS_DIR'], wid, filename)
    if file_exists(path):
        filename = name + '_' + str(int(time())) + ext
        path = path_join(dirname(path), filename)                               # TODO: repeat line 170
    file.save(path)                                                             # TODO: handle failed to save
    file.close()
    up = Upload(filename=filename, size=getsize(path),
                type=filetype, created=datetime.fromtimestamp(getctime(path)),       # TODO: use datetime.now()
                worker_id=wid)
    db.session.add(up)
    db.session.commit()
    return '', 201, {'Location': url_for('get_single_upload', wid=wid, uid=up.id)}


@app.route('/api/workers/<wid>/uploads/<int:uid>', methods=['GET'])
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


@app.route('/api/workers/<wid>/uploads/<int:uid>', methods=['DELETE'])  # TODO: add own parameter decoder?
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


@app.route('/api/workers/<wid>/uploads/<int:uid>/info', methods=['GET'])
def get_upload_info(wid, uid):
    if uid == 0:
        return jsonify(Upload.query.filter_by(worker_id=wid).all()), 200
    up = Upload.query.filter_by(worker_id=wid, id=uid)
    return (jsonify(up), 200) if up else ('', 404)


# REPORTS --------------------------------------------------------------------
@app.route('/api/workers/<wid>/jobs/<int:jid>/report', methods=['POST'])
def create_report(wid, jid):
    if request.mimetype != 'text/plain' or request.content_length > app.config['REPORT_LEN']:
        return '', 422
    if not obj_exists(and_(Job.worker_id == wid, Job.id == jid)):
        return '', 404
    rep = JobReport(job_id=jid, report=request.get_data(cache=False))
    try:
        db.session.add(rep)
        Job.query.filter_by(id=jid).update({'is_done': True}, synchronize_session=False)
        db.session.commit()
    except IntegrityError:  # report already exists, means Job.is_done == True
        db.session.rollback()
        return '', 409
    return '', 201, {'Location': url_for('get_report', wid=wid, jid=jid)}


@app.route('/api/workers/<wid>/jobs/<int:jid>/report', methods=['GET'])
def get_report(wid, jid):
    if not obj_exists(and_(Job.worker_id == wid, Job.id == jid)):
        return '', 404
    rep = JobReport.query.get(jid)
    return (rep.report, 200, {'Content-type': 'text/plain; charset=utf-8'}) if rep else ('', 404, {})


@app.route('/api/workers/<wid>/jobs/<int:jid>/report', methods=['DELETE'])
def delete_report(wid, jid):
    if not obj_exists(and_(Job.worker_id == wid, Job.id == jid)):
        return '', 404
    if JobReport.query.filter_by(job_id=jid).delete(synchronize_session=False) == 0:
        return '', 404
    db.session.commit()
    return '', 200
