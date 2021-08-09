from flask import Flask
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
app.config.from_pyfile('../config.py')  # load config in run.py

db = SQLAlchemy(app)
# TODO: in sqlite3 enable foreign keys with pragma
