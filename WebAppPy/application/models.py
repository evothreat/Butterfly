from datetime import datetime
from base64 import urlsafe_b64encode
from dataclasses import dataclass
from uuid import uuid4

from sqlalchemy import exists
from flask_login import UserMixin
from application import db


def obj_exists(cond):
    return db.session.query(exists().where(cond)).scalar()


def base64_uuid():
    return urlsafe_b64encode(uuid4().bytes).rstrip(b'=').decode('ascii')


@dataclass
class JobReport(db.Model):
    job_id: int
    report: str

    job_id = db.Column(db.Integer, db.ForeignKey('job.id'), primary_key=True, nullable=False)
    report = db.Column(db.Text)


@dataclass
class Job(db.Model):
    id: int
    todo: str
    is_done: bool
    created: datetime
    worker_id: str

    id = db.Column(db.Integer, primary_key=True)
    todo = db.Column(db.String(250))
    is_done = db.Column(db.Boolean, default=False)
    created = db.Column(db.DateTime, default=datetime.now)
    worker_id = db.Column(db.String(22), db.ForeignKey('worker.id'), nullable=False)

    @staticmethod
    def from_dict(d):
        return Job(todo=d['todo'])


@dataclass
class HardwareInfo(db.Model):
    gpu: str
    cpu: str
    ram: str
    worker_id: str

    gpu = db.Column(db.String(48))
    cpu = db.Column(db.String(64))
    ram = db.Column(db.String(10))                                                  # TODO: maybe use integer?
    worker_id = db.Column(db.String(22), db.ForeignKey('worker.id'), primary_key=True, nullable=False)

    @staticmethod
    def from_dict(d):
        return HardwareInfo(gpu=d['gpu'], cpu=d['cpu'], ram=d['ram'])


@dataclass
class Worker(db.Model):
    id: str
    hostname: str
    country: str
    os: str
    ip_addr: str
    is_admin: bool
    boost: bool
    last_seen: datetime

    id = db.Column(db.String(22), primary_key=True)
    hostname = db.Column(db.String(30))
    ip_addr = db.Column(db.String(15))
    country = db.Column(db.String(15))
    os = db.Column(db.String(15))
    is_admin = db.Column(db.Boolean)
    boost = db.Column(db.Boolean, default=False)
    last_seen = db.Column(db.DateTime, default=datetime.now)

    @staticmethod
    def from_dict(d):
        return Worker(id=d['id'], hostname=d['hostname'], country=d['country'], ip_addr=d['ip_addr'],
                      os=d['os'], is_admin=d['is_admin'])


@dataclass
class Upload(db.Model):
    id: int
    filename: str
    type: str
    size: int
    created: datetime
    worker_id: str
    # TODO: add content url?

    id = db.Column(db.Integer, primary_key=True)
    filename = db.Column(db.String(64))
    type = db.Column(db.String(16))
    size = db.Column(db.BigInteger)
    created = db.Column(db.DateTime)
    worker_id = db.Column(db.String(22), db.ForeignKey('worker.id'), nullable=False)


class Admin(UserMixin, db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(30), unique=True)
    password = db.Column(db.String(64))
