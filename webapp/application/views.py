from application import *
from application.models import *
from flask import jsonify


@app.route('/')
def hello_world():  # put application's code here
    return 'Hello World!'


@app.route('/api/v1/workers', methods=['GET'])
def get_workers():
    return jsonify(Worker.query.all()), 200


@app.route('/api/v1/workers', methods=['POST'])
def create_worker():
    pass


@app.route('/api/v1/workers/<int:wid>', methods=['GET'])
def get_single_worker(wid):
    w = Worker.query.get(wid)
    return ('', 404) if not w else (jsonify(w), 200)


@app.route('/api/v1/workers/<int:wid>', methods=['DELETE'])
def delete_worker(wid):
    Worker.query.filter_by(id=wid).delete()
    db.session.commit()
    return '', 200


@app.route('/api/v1/workers/<int:wid>/jobs', methods=['GET'])
def get_jobs(wid):
    w = Worker.query.get(wid)
    return ('', 404) if not w else (jsonify(w.jobs), 200)


@app.route('/api/v1/workers/<int:wid>/jobs', methods=['POST'])
def create_job(wid):
    pass
