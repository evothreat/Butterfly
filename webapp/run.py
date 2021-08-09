from application import *
import application.models as models
import application.api
import application.control


def add_test_data():
    j = models.Job(args='upload test.exe', worker_id=1)
    w = models.Worker(hostname='Predator', os='Windows 10', country='Germany', ip_addr='127.0.0.1',
                      mac_addr='fe80::9ec8:fcff:fe3e:daba')
    w2 = models.Worker(hostname='Helios 300', os='Windows 7', country='England', ip_addr='127.2.4.1',
                       mac_addr='fe80::9ec8:hgff:fe3e:daba')
    w3 = models.Worker(hostname='Acer Nexus', os='Windows 8', country='USA', ip_addr='127.11.55.1',
                       mac_addr='fe80::9ec8:vdsa:fe3e:daba')
    w4 = models.Worker(hostname='Predator', os='Windows 10', country='Russia', ip_addr='127.66.33.1',
                       mac_addr='fe80::9ec8:abcd:fe3e:daba')
    db.session.add(w)
    db.session.add(w2)
    db.session.add(w3)
    db.session.add(w4)
    db.session.add(j)


if __name__ == '__main__':
    # db.drop_all()
    # db.create_all()
    # add_test_data()
    # db.session.commit()
    app.run()
