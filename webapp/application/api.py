from sqlalchemy.exc import IntegrityError

from application import *
from application.models import *
from flask import jsonify, request


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
    except KeyError:
        return '', 422
    except IntegrityError:
        return '', 409
    return jsonify(w), 200


@app.route('/api/v1/workers/<int:wid>', methods=['GET'])
def get_single_worker(wid):
    w = Worker.query.get(wid)
    return ('', 404) if not w else (jsonify(w), 200)


@app.route('/api/v1/workers/<int:wid>', methods=['DELETE'])
def delete_worker(wid):
    if Worker.query.filter_by(id=wid).delete() == 0:
        return '', 404
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<int:wid>/jobs', methods=['GET'])
def get_jobs(wid):
    if not Worker.query.get(wid):
        return '', 404
    args = request.args.get('is_done')
    is_done = bool(args and args.lower() == 'true')
    jobs = Job.query.filter_by(worker_id=wid, is_done=is_done).all()
    return jsonify(jobs), 200


@app.route('/api/v1/workers/<int:wid>/jobs', methods=['POST'])
def create_job(wid):
    if not request.is_json:
        return '', 400
    if not Worker.query.get(wid):  # TODO: remove after enabling foreign keys
        return '', 404
    try:
        j = Job.from_dict(request.json)
        j.worker_id = wid
    except KeyError:
        return '', 422
    db.session.add(j)
    db.session.commit()
    return jsonify(j), 200
