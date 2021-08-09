from application import *
from application.models import *
from flask import render_template


# TODO: authentication required
@app.route('/worker-list', methods=['GET', 'POST'])
def list_workers():
    return render_template('worker-list.html', title='Worker-List', workers=Worker.query.all())
