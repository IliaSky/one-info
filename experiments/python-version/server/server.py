from flask import Flask
from pymongo import MongoClient
from os import getcwd

import gmail

client = MongoClient()
db = client.test_database
posts = db.posts
app = Flask(__name__, static_folder="../client/", static_url_path="")


@app.route('/')
def root():
    return app.send_static_file('index.html')


@app.route("/db/add/<sth>")
def add(sth):
    posts.insert({"content": sth})
    return "Added " + sth


@app.route("/db/all")
def list_all():
    return str(list(posts.find()))


# @app.route("/api/folders")
# def folders():
#     return gmail.api_folders()


@app.route("/api/folders")
def folders():
    return gmail.api(gmail.directories)


@app.route("/api/folders/<sth>")
def folder(sth):
    try:
        return gmail.api(gmail.folder, sth)
    except Exception as e:
        print(e)
        return str(e)


if __name__ == "__main__":
    app.run()