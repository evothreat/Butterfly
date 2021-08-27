from datetime import date

from flask import Flask
from flask.json import JSONEncoder
from flask_sqlalchemy import SQLAlchemy
from flask_login import LoginManager


class JsonEncoder(JSONEncoder):         # TODO: export this class?
    def default(self, obj):
        if isinstance(obj, date):
            return obj.strftime('%a, %e %b %Y %H:%M:%S')
        return super().default(obj)


app = Flask(__name__)
app.config.from_pyfile('../config.py')  # load config in run.py
app.json_encoder = JsonEncoder

login_manager = LoginManager(app)
login_manager.login_view = '/login'  # TODO: use url_for

db = SQLAlchemy(app)
# TODO: in sqlite3 enable foreign keys with pragma
