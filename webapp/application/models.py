from application import db
import datetime
from dataclasses import dataclass
from sqlalchemy.sql import func
from flask_login import UserMixin


@dataclass
class Job(db.Model):
    id: int
    args: str
    is_done: bool
    timestamp: datetime.datetime
    worker_id: int

    id = db.Column(db.Integer, primary_key=True)
    args = db.Column(db.String(250))
    is_done = db.Column(db.Boolean(), default=False)
    timestamp = db.Column(db.DateTime(), default=func.now())
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'), nullable=False)

    @staticmethod
    def from_dict(d):
        return Job(args=d['args'])


@dataclass
class ResourceInfo(db.Model):
    id: int
    gpu: str
    cpu: str
    ram: str
    worker_id: int

    id = db.Column(db.Integer, primary_key=True)
    gpu = db.Column(db.String(48))
    cpu = db.Column(db.String(64))
    ram = db.Column(db.String(10))  # TODO: maybe use integer?
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'), nullable=False)

    @staticmethod
    def from_dict(d):
        return ResourceInfo(gpu=d['gpu'], cpu=d['cpu'], ram=d['ram'])


@dataclass
class Worker(db.Model):
    id: int
    hostname: str
    os: str
    country: str
    ip_addr: str
    mac_addr: str
    last_seen: datetime.datetime
    # resource_info: ResourceInfo
    # jobs: Job

    id = db.Column(db.Integer, primary_key=True)
    hostname = db.Column(db.String(30))
    os = db.Column(db.String(15))
    country = db.Column(db.String(15))
    ip_addr = db.Column(db.String(15))
    mac_addr = db.Column(db.String(17), unique=True)  # set as primary key?
    last_seen = db.Column(db.DateTime, default=func.now())

    # resource_info = db.relationship("ResourceInfo", uselist=False, cascade="all, delete-orphan")
    # jobs = db.relationship("Job", cascade="all, delete-orphan")

    @staticmethod
    def from_dict(d):
        return Worker(hostname=d['hostname'], os=d['os'], country=d['country'],
                      ip_addr=d['ip_addr'], mac_addr=d['mac_addr'])


@dataclass
class Upload(db.Model):
    id: int
    filename: str
    type: str
    size: int
    created: datetime.datetime
    worker_id: int
    # TODO: add content url?

    id = db.Column(db.Integer, primary_key=True)
    filename = db.Column(db.String(64))
    type = db.Column(db.String(16))
    size = db.Column(db.BigInteger())
    created = db.Column(db.DateTime)
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'), nullable=False)


class Admin(UserMixin, db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(30), unique=True)
    password = db.Column(db.String(64))
