from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_login import LoginManager

app = Flask(__name__)
app.config.from_pyfile('../config.py')  # load config in run.py

login_manager = LoginManager(app)
login_manager.login_view = '/login'     # TODO: use url_for

db = SQLAlchemy(app)
# TODO: in sqlite3 enable foreign keys with pragma
