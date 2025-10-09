package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Item struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	SmallPrice  int           `bson:"s_price"`
	BigPrice    int           `bson:"b_price"`
	Description string        `bson:"description"`
	Ingredients []string      `bson:"ingredients"`
	Type        string        `bson:"type"`
	Class       string        `bson:"class"`
	ImageLink   string
	Url         string
	SmallPriceF string
	BigPriceF   string
	IdHex       string
}

const URI = "mongodb://localhost:27017/"

func main() {

	log.SetFlags(log.Lshortfile | log.Ltime)

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

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cur, err := collection.Find(ctx, bson.M{"type": "boloCaseiro"})
		if err != nil {

			log.Println(err)
			c.Status(500)
			return
		}

		var total []Item
		for cur.Next(ctx) {
			var v Item
			if err := cur.Decode(&v); err != nil {
				log.Println(err)
				c.Status(500)
				return
			}

			v.Url = fmt.Sprintf("https://servidordomal.fun/produto/%s", v.ID.Hex())
			v.ImageLink = fmt.Sprintf("https://servidordomal.fun/static/imgs/%s.jpg", v.ID.Hex())

			fmt.Println(v.Name)
			total = append(total, v)
		}

		c.HTML(http.StatusOK, "ginTemplateFormat.html", gin.H{"bolosCaseiros": total})
	})

	r.GET("/produto/:id", func(c *gin.Context) {

		id, err := bson.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var v Item
		if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&v); err != nil {
			log.Println(err)
			c.Status(500)
			return
		}

		v.ImageLink = fmt.Sprintf("https://servidordomal.fun/static/imgs/%s.jpg", c.Param("id"))

		v.SmallPriceF = fmt.Sprintf("R$ %.2f", float64(v.SmallPrice)/100)
		v.BigPriceF = fmt.Sprintf("R$ %.2f", float64(v.BigPrice)/100)
		c.HTML(200, "ginDetailsTemplate.html", gin.H{"Item": v})
	})

	r.GET("/admin/add", func(c *gin.Context) {
		c.HTML(200, "Add.html", gin.H{})
	})
	r.GET("/admin/add?success=true", func(c *gin.Context) {
		c.HTML(200, "Success.html", gin.H{})
	})

	r.GET("/admin/edit/:id", func(c *gin.Context) {
		idHex := c.Param("id")

		id, err := bson.ObjectIDFromHex(idHex)
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var item Item

		item.IdHex = idHex

		if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
			log.Println(err)
			c.Status(400)
			return
		}

		jsonB, err := json.Marshal(item.Ingredients)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return
		}

		c.HTML(200, "ginEdit.html", gin.H{
			"Item":            item,
			"JSONIngredients": template.JS(string(jsonB)),
		})
	})

	r.POST("/admin/add", func(c *gin.Context) {
		var item Item
		var err error

		item.Name = c.PostForm("name")
		item.SmallPrice, err = strconv.Atoi(c.PostForm("s_price"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		item.BigPrice, err = strconv.Atoi(c.PostForm("b_price"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		item.Description = c.PostForm("description")
		if len(c.PostForm("ingredients")) <= 0 {
			c.JSON(400, gin.H{"erro": "o produto precisa ter ingredientes"})
			return
		}
		if err := json.Unmarshal([]byte(c.PostForm("ingredients")), &item.Ingredients); err != nil {
			log.Println(err)
			c.Status(400)
			return
		}

		item.Type = c.PostForm("type")

		file, err := c.FormFile("image")
		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"erro": "é necessário uma imagem para adicionar um produto"})
			return
		}

		if !strings.HasSuffix(file.Filename, ".jpg") {
			c.JSON(400, gin.H{"erro": "somente .jpg permitido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := collection.InsertOne(ctx, item)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return
		}

		id, ok := res.InsertedID.(bson.ObjectID)
		if !ok {
			log.Println("failed to get objectid")
			c.Status(500)
			return
		}

		if err := c.SaveUploadedFile(file, fmt.Sprintf("./static/imgs/%s.jpg", id.Hex()), 0777); err != nil {
			log.Println(err)
			c.Status(500)
			return
		}

		c.Redirect(303, "/admin/add?success=true")
	})

	r.POST("/admin/edit", func(c *gin.Context) {
		var item Item
		var err error

		item.ID, err = bson.ObjectIDFromHex(c.PostForm("id"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		item.Name = c.PostForm("name")
		item.SmallPrice, err = strconv.Atoi(c.PostForm("s_price"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		item.BigPrice, err = strconv.Atoi(c.PostForm("b_price"))
		if err != nil {
			log.Println(err)
			c.Status(400)
			return
		}
		item.Description = c.PostForm("description")
		if len(c.PostForm("ingredients")) <= 0 {
			c.JSON(400, gin.H{"erro": "o produto precisa ter ingredientes"})
			return
		}
		if err := json.Unmarshal([]byte(c.PostForm("ingredients")), &item.Ingredients); err != nil {
			log.Println(err)
			c.Status(400)
			return
		}

		item.Type = c.PostForm("type")

		file, err := c.FormFile("image")
		if err == nil {
			if !strings.HasSuffix(file.Filename, ".jpg") {
				c.JSON(400, gin.H{"erro": "somente .jpg permitido"})
				return
			}
			if err := c.SaveUploadedFile(file, fmt.Sprintf("./static/imgs/%s.jpg", item.ID.Hex()), 0777); err != nil {
				log.Println(err)
				c.Status(500)
				return
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection.ReplaceOne(ctx, bson.M{"_id": item.ID}, item)

		c.Redirect(303, "/admin")
	})

	if err := r.RunTLS(":443", "/etc/letsencrypt/live/servidordomal.fun/fullchain.pem", "/etc/letsencrypt/live/servidordomal.fun/privkey.pem"); err != nil {
		panic(err)
	}

}
