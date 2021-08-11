from application import *
import application.models as models
import application.api
import application.control
from werkzeug.security import generate_password_hash
from time import sleep


def add_test_data():
    j = models.Job(args='upload test.exe', worker_id=1)

    ri = models.ResourceInfo(cpu='Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz 2.21 GHz',
                             gpu='NVIDIA GeForce GTX 1060',
                             ram='8,0 GB', worker_id=1)
    ri2 = models.ResourceInfo(cpu='Intel Core i9-10900K Processor',
                              gpu='AMD Radeon RX 5600 XT',
                              ram='4,0 GB', worker_id=2)
    ri3 = models.ResourceInfo(cpu='AMD Ryzen 9 5900X',
                              gpu='Nvidia GeForce RTX 3080',
                              ram='12,0 GB', worker_id=3)
    ri4 = models.ResourceInfo(cpu='Intel Core i9-10980XE Extreme Edition Processor',
                              gpu='Intel(R) UHD Graphics 630',
                              ram='32,0 GB', worker_id=4)

    w = models.Worker(hostname='Predator', os='Windows 10', country='Germany', ip_addr='127.0.0.1',
                      mac_addr='fe80::9ec8:fcff:fe3e:daba')
    w2 = models.Worker(hostname='Helios 300', os='Windows 7', country='England', ip_addr='127.2.4.1',
                       mac_addr='fe80::9ec8:hgff:fe3e:daba')
    w3 = models.Worker(hostname='Acer Nexus', os='Windows 8', country='USA', ip_addr='127.11.55.1',
                       mac_addr='fe80::9ec8:vdsa:fe3e:daba')
    w4 = models.Worker(hostname='Predator', os='Windows 10', country='Russia', ip_addr='127.66.33.1',
                       mac_addr='fe80::9ec8:abcd:fe3e:daba')
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

    db.session.add(j)
    db.session.add(ri)
    db.session.add(ri2)
    db.session.add(ri3)
    db.session.add(ri4)

    admin = models.Admin()
    admin.username = 'Adam'
    admin.password = generate_password_hash('12345')
    db.session.add(admin)


if __name__ == '__main__':
    # db.drop_all()
    # db.create_all()
    # add_test_data()
    # db.session.commit()
    app.run()
