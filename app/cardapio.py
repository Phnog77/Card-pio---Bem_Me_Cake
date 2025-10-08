#!/usr/bin/env python
# -*- coding: utf-8 -*-
from flask import Flask, render_template, jsonify
from flask_cors import CORS
from pymongo import MongoClient

app = Flask(__name__)
CORS(app)

client = MongoClient("mongodb://localhost:27017/")
db = client["BemMeCake"]
produtos = db["items"]

@app.route("/cardapio/")
def index():
    bolos_caseiros = db.session.query(BolosCaseiros).all()  # Exemplo de consulta
    bolos_de_festa = db.session.query(BolosDeFesta).all()
    doces = db.session.query(Doces).all()
    recheios = db.session.query(Recheios).all()
    bolos_de_pote = db.session.query(BolosDePote).all()
    novidades = db.session.query(Novidades).all()
    bebidas = db.session.query(Bebidas).all()
    
    return render_template('index.html', 
                           bolos_caseiros=bolos_caseiros, 
                           bolos_de_festa=bolos_de_festa,
                           doces=doces, 
                           recheios=recheios, 
                           bolos_de_pote=bolos_de_pote, 
                           novidades=novidades, 
                           bebidas=bebidas)


if __name__ == "__main__":
    app.run(port=31520, debug=True, host="127.0.0.1")
