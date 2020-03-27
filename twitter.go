package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func authTwitter(consumerKey string, consumerSecret string, accessToken string, accessSecret string) *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	twitterClient := twitter.NewClient(httpClient)

	return twitterClient
}

func getProfile(twitterClient *twitter.Client, myProfile *userModel) error {
	user, _, err := twitterClient.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		return err
	}

	followers, _, err := twitterClient.Followers.IDs(&twitter.FollowerIDParams{})
	if err != nil {
		return err
	}

	ids := make([]string, len(followers.IDs))
	for i, id := range followers.IDs {
		ids[i] = string(id)
	}

	myProfile = &userModel{
		Name:           user.ScreenName,
		UserProvider:   "twitter",
		UserID:         user.IDStr,
		Friends:        ids,
		ImageURI:       user.ProfileImageURL,
		NumOfFollowers: user.FollowersCount,
	}
	return nil
}
