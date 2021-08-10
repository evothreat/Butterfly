from application import *
from application.models import *
from flask import render_template, request, redirect
from werkzeug.security import check_password_hash
from flask_login import login_user, login_required, current_user


@login_manager.user_loader
def load_user(uid):
    return Admin.query.get(uid)


@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'GET':
        return render_template('login.html')
    name = request.form.get('username')
    passwd = request.form.get('password')
    # TODO: check if name and passwd empty
    remember = bool(request.form.get('remember'))
    admin = Admin.query.filter_by(username=name).first()
    if admin and check_password_hash(admin.password, passwd):
        # flash message
        login_user(admin, remember)
        return redirect('/worker-list')  # TODO: use url_for
    return redirect('/login')


# TODO: authentication required
@app.route('/worker-list', methods=['GET', 'POST'])
@login_required
def list_workers():
    return render_template('worker-list.html', title='Worker-List', workers=Worker.query.all())
