import application.views
import application.models

if __name__ == '__main__':
    application.db.create_all()
    application.db.session.commit()
    application.app.run()
