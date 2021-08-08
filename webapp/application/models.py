from application import db


class ResourceInfo(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    gpu = db.Column(db.String(50))
    cpu = db.Column(db.String(50))
    ram = db.Column(db.String(10))                                          # maybe use integer?
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'))


class Worker(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    hostname = db.Column(db.String(30))
    os = db.Column(db.String(30))
    country = db.Column(db.String(30))
    ip_addr = db.Column(db.String(15))
    mac_addr = db.Column(db.String(17), unique=True)                                     # set as primary key?
    last_seen = db.Column(db.DateTime())
    resource_info = db.relationship("ResourceInfo", uselist=False)
    jobs = db.relationship("Job")


class Job(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    args = db.Column(db.String(250))
    completed = db.Column(db.Boolean())
    timestamp = db.Column(db.DateTime())
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'))
