from application import *
from application.models import *
from flask import render_template, request


@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'GET':
        return render_template('login.html')


# TODO: authentication required
@app.route('/worker-list', methods=['GET', 'POST'])
def list_workers():
    return render_template('worker-list.html', title='Worker-List', workers=Worker.query.all())
