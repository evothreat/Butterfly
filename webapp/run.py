from application import *
import application.models as models
import application.views as views


def add_test_data():
    j = models.Job(args='upload test.exe', worker_id=1)
    w = models.Worker(hostname='Predator', os='Windows 10', country='Germany', ip_addr='127.0.0.1',
                      mac_addr='fe80::9ec8:fcff:fe3e:daba')
    db.session.add(w)
    db.session.add(j)


if __name__ == '__main__':
    db.drop_all()
    db.create_all()
    add_test_data()
    db.session.commit()
    app.run()
