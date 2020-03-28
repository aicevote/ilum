package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dinever/golf"
	"github.com/jakehl/goid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database = nil
var consumerKey, consumerSecret = "", ""

type authRequest struct {
	AccessToken  string
	AccessSecret string
}

type authResponce struct {
	SessionID string
}

func mainHandler(ctx *golf.Context) {
	data := authRequest{}
	decoder := json.NewDecoder(ctx.Request.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		ctx.Abort(400)
		return
	}

	myProfile := userModel{}
	twitterClient := authTwitter(consumerKey, consumerSecret, data.AccessToken, data.AccessSecret)
	err = getProfile(twitterClient, &myProfile)
	if err != nil {
		fmt.Println(err)
		ctx.Abort(400)
		return
	}

	sessionID := goid.NewV4UUID().String()
	err = saveUserData(&myProfile, data.AccessToken, data.AccessSecret, sessionID)
	if err != nil {
		fmt.Println(err)
		ctx.Abort(500)
		return
	}

	ctx.JSON(authResponce{SessionID: sessionID})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")

	uri := os.Getenv("DB_URI")
	client, err := connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("glacierapi")

	app := golf.New()
	app.Post("/ilum", mainHandler)
	app.Run(":9000")

	fmt.Println("Hello World!")
}
