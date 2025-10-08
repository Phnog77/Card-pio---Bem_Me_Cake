#!/usr/bin/env python
# -*- coding: utf-8 -*-
# |
# Imports
# | (Flask)
from flask import Flask, request, jsonify, send_from_directory, make_response, render_template, redirect
from flask_cors import CORS
from flask import send_file
# | (Others) 
import os
import jwt
import time
import json
import uuid
import flask
import socket
import random
import bcrypt
import shutil
import requests
import threading, pytz
from random import randint
from threading import Timer
from threading import Thread
from cryptography.fernet import Fernet
from datetime import datetime, timedelta
from werkzeug.utils import secure_filename
from jwt.exceptions import ExpiredSignatureError, InvalidTokenError
# |
# |
app = Flask(__name__)
CORS(app)
# |
@app.route("/aps/carregar")
def carregar():
    
# |
# |
# |
# Start API
if __name__ == '__main__':
    app.run(port=31523, debug=True, host="127.0.0.1")