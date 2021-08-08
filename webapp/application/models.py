import datetime
from application import db
from dataclasses import dataclass
from sqlalchemy.sql import func


@dataclass
class Job(db.Model):
    id: int
    args: str
    completed: bool
    timestamp: datetime.datetime
    worker_id: int

    id = db.Column(db.Integer, primary_key=True)
    args = db.Column(db.String(250))
    completed = db.Column(db.Boolean(), default=False)
    timestamp = db.Column(db.DateTime(), default=func.now())
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'))


@dataclass
class ResourceInfo(db.Model):
    id: int
    gpu: str
    cpu: str
    ram: str
    worker_id: int

    id = db.Column(db.Integer, primary_key=True)
    gpu = db.Column(db.String(50))
    cpu = db.Column(db.String(50))
    ram = db.Column(db.String(10))  # maybe use integer?
    worker_id = db.Column(db.Integer, db.ForeignKey('worker.id'))


@dataclass
class Worker(db.Model):
    id: int
    hostname: str
    os: str
    country: str
    ip_addr: str
    mac_addr: str
    last_seen: datetime.datetime
    resource_info: ResourceInfo
    jobs: Job

    id = db.Column(db.Integer, primary_key=True)
    hostname = db.Column(db.String(30))
    os = db.Column(db.String(30))
    country = db.Column(db.String(30))
    ip_addr = db.Column(db.String(15))
    mac_addr = db.Column(db.String(17), unique=True)  # set as primary key?
    last_seen = db.Column(db.DateTime(), default=func.now())
    resource_info = db.relationship("ResourceInfo", uselist=False, cascade="all, delete-orphan")
    jobs = db.relationship("Job", cascade="all, delete-orphan")
