package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Item struct {
	Name        string   `bson:"name"`
	SmallPrice  int      `bson:"s_price"`
	BigPrice    int      `bson:"b_price"`
	Description string   `bson:"description"`
	Ingredients []string `bson:"ingredients"`
	Type        string   `bson:"type"`
	Image       string   `bson:"image"`
	Class       string   `bson:"class"`
	ImageLink   string
	Url         string
}

const URI = "mongodb://localhost:27017/"

func main() {
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("BemMeCake")
	collection := db.Collection("items")

	r.Static("/static", "../static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {

		cur, err := collection.Find(ctx, bson.M{"type": "bolo", "class": "caseira"})
		if err != nil {
			log.Fatal(err)
		}

		var total []Item
		for cur.Next(ctx) {
			var v Item
			if err := cur.Decode(&v); err != nil {
				log.Fatal(err)
			}

			total = append(total, v)
		}

		c.HTML(http.StatusOK, "ginTemplateFormat.html", gin.H{"bolosCaseiros": total})
	})

	if err := r.RunTLS(":443", "/etc/letsencrypt/live/servidordomal.fun/fullchain.pem", "/etc/letsencrypt/live/servidordomal.fun/privkey.pem"); err != nil {
		panic(err)
	}

}
