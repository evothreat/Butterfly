from application import *
from application.models import *
from flask import render_template, request, redirect
from werkzeug.security import check_password_hash
from flask_login import login_user, login_required, logout_user


@login_manager.user_loader
def load_user(uid):
    return Admin.query.get(uid)


@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'GET':
        return render_template('login.html')
    try:
        name = request.form['username']
        passwd = request.form['password']
    except KeyError:
        return redirect('/login')
    remember = bool(request.form.get('remember'))
    admin = Admin.query.filter_by(username=name).first()
    if admin and check_password_hash(admin.password, passwd):
        login_user(admin, remember)
        return redirect('/workers')  							# TODO: use url_for
    return redirect('/login')


@app.route('/logout', methods=['POST'])
def logout():
    logout_user()
    return redirect('/login')


@app.route('/workers', methods=['GET', 'POST'])     			# TODO: change to index?
@login_required
def list_workers():
    return render_template('workers.html', title='Workers')


@app.route('/workers/<wid>')
@login_required
def interact(wid):
    w = Worker.query.get(wid)
    if not w:
        return redirect('/workers', code=302)
    return render_template('interact.html', title='Interaction')
