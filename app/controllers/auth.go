package controllers

import (
	"encoding/json"
	"github.com/mrjones/oauth"
	"github.com/revel/revel"
	"twitterApp/app/models"
)

type Auth struct {
	*revel.Controller
}

// oauthのコンシューマ設定
var twitter = oauth.NewConsumer(
	"consumerkey",
	"consumersecret"
	oauth.ServiceProvider{
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)

func (c Auth) Index() revel.Result {
	// Twitterから情報を取得する
	user := getUser()

	// callback URLを設定する
	requestToken, url, err := twitter.GetRequestTokenAndUrl("http://localhost:9000/auth/callback")
	if err == nil {
		// ユーザ情報セット
		user.RequestToken = requestToken
		// oauth_verifierを取得
		return c.Redirect(url)
	}
	revel.ERROR.Println("リクエストトークン取得できませんでした！！", err)

	return c.Render()
}

func (c Auth) Callback(oauth_verifier string) revel.Result {
	// セットしたユーザ情報取得
	user := getUser()

	// access_tokenを獲得
	accessToken, err := twitter.AuthorizeToken(user.RequestToken, oauth_verifier)
	if err == nil {
		// ユーザ情報を取得する
		resp, _ := twitter.Get(
			"https://api.twitter.com/1.1/account/verify_credentials.json",
			map[string]string{},
			accessToken,
		)
		defer resp.Body.Close()
		account := struct {
			Name            string `json:"name"`
			ProfileImageURL string `json:"profile_image_url"`
		}{}
		_ = json.NewDecoder(resp.Body).Decode(&account)

		// 表示用情報をセット
		setUserData(account.Name, account.ProfileImageURL)
	} else {
		// 失敗
		revel.ERROR.Println("取得失敗！！", err)
	}

	return c.Redirect(Show.Index)
}

func (c Auth) Show() revel.Result {
	// ユーザ情報取得
	user := getShowUser("kaisou_test")
	return c.Render(user)
}

// Twitterユーザ情報
func getUser() *models.User {
	return models.FindOrCreate("kaisou")
}

// 表示用ユーザ情報セット
func setUserData(name, imgURL string) {
	models.CreateShowUser(name, imgURL)
}
