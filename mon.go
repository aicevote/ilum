package main

import (
	"context"
	"time"

	"github.com/jakehl/goid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userModel struct {
	Name           string
	UserProvider   string
	UserID         string
	Friends        []string
	ImageURI       string
	NumOfFollowers int
}

type sessionModel struct {
	UserProvider       string
	UserID             string
	AccessToken        string
	RefreshToken       string
	SessionID          string
	SessionIDExpire    int
	SessionToken       string
	SessionTokenExpire int
}

func connect(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func saveUserData(user *userModel, accessToken string, refreshToken string, sessionID string) error {
	users := db.Collection("users")
	sessions := db.Collection("sessions")

	filter := bson.D{{"userProvider", user.UserProvider}, {"userID", user.UserID}}

	update := bson.D{
		{"$set", bson.D{
			{"name", user.Name},
			{"friends", user.Friends},
			{"imageURI", user.ImageURI},
			{"numOfFollowers", user.NumOfFollowers},
		}},
	}

	upsert := true
	_, err := users.UpdateOne(context.TODO(),
		filter, update, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		return err
	}

	oneHour := 60 * 60 * 1000
	oneDay := oneHour * 24
	oneMonth := oneDay * 31

	now := int(time.Now().UnixNano() / int64(time.Millisecond))

	session := sessionModel{
		UserProvider:       user.UserProvider,
		UserID:             user.UserID,
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		SessionID:          sessionID,
		SessionIDExpire:    now + oneMonth,
		SessionToken:       goid.NewV4UUID().String(),
		SessionTokenExpire: now + oneHour,
	}

	_, err = sessions.InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}

	return nil
}
