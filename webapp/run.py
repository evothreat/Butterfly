from application import *
import application.models as models
import application.views as views


def add_test_data():
    job = models.Job(args='upload test.exe')
    worker = models.Worker(hostname='Predator', os='Windows 10', country='Germany', ip_addr='127.0.0.1',
                           mac_addr='fe80::9ec8:fcff:fe3e:daba')
    worker.jobs.append(job)
    db.session.add(worker)


if __name__ == '__main__':
    db.create_all()
    add_test_data()
    db.session.commit()
    app.run()
