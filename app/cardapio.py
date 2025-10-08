#!/usr/bin/env python
# -*- coding: utf-8 -*-
from flask import Flask, render_template, jsonify
from flask_cors import CORS
from pymongo import MongoClient

app = Flask(__name__)
CORS(app)

client = MongoClient("mongodb://" + "admin" + open("senha.txt").read() + ":admin@localhost:27017/")
db = client["BemMeCake"]
produtos = db["items"]

@app.route("/cardapio/")
def carregar():
    bolos_caseiros = list(produtos.find({"type": "bolo"}))
    doces = list(produtos.find({"type": "doce"}))
    bolos_de_festa = list(produtos.find({"type": "festa"}))
    recheios = list(produtos.find({"type": "recheio"}))
    bolos_de_pote = list(produtos.find({"type": "pote"}))
    novidades = list(produtos.find({"type": "novidade"}))
    bebidas = list(produtos.find({"type": "bebida"}))

    return render_template(
        "index",
        bolos_caseiros=bolos_caseiros,
        doces=doces,
        bolos_de_festa=bolos_de_festa,
        recheios=recheios,
        bolos_de_pote=bolos_de_pote,
        novidades=novidades,
        bebidas=bebidas,
    )

if __name__ == "__main__":
    app.run(port=31520, debug=True, host="127.0.0.1")
