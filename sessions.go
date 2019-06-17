// sessions.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	store        = sessions.NewCookieStore([]byte(os.Getenv("SESS_SECRET")))
	oauth1Config = oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
)

func main() {
	oauth1Config.CallbackURL = os.Getenv("CALLBACK_URL")
	oauth1Config.Endpoint.RequestTokenURL = "https://api.twitter.com/oauth/request_token"
	oauth1Config.Endpoint.AuthorizeURL = "https://api.twitter.com/oauth/authenticate"
	oauth1Config.Endpoint.AccessTokenURL = "https://api.twitter.com/oauth/access_token"

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello!") })
	http.HandleFunc("/oauth/request_token", requestToken)
	http.HandleFunc("/oauth/callback", callback)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func requestToken(w http.ResponseWriter, r *http.Request) {
	session, e := store.Get(r, "cookie-name")
	if e != nil {
		panic(e)
	}
	session.Values["consumer_key"] = os.Getenv("CONSUMER_KEY")
	session.Values["consumer_secret"] = os.Getenv("CONSUMER_SECRET")
	session.Values["state"] = r.URL.Query().Get("state")
	session.Values["client_id"] = r.URL.Query().Get("client_id")
	session.Values["redirect_uri"] = r.URL.Query().Get("redirect_uri")

	requestToken, requestTokenSecret, e := oauth1Config.RequestToken()
	if e != nil {
		panic(e)
	}
	session.Values["request_token"] = requestToken
	session.Values["request_token_secret"] = requestTokenSecret

	e = session.Save(r, w)
	if e != nil {
		panic(e)
	}
	http.Redirect(w, r, "https://twitter.com/oauth/authenticate?oauth_token="+requestToken, http.StatusFound)
}

func callback(w http.ResponseWriter, r *http.Request) {
	session, e := store.Get(r, "cookie-name")
	if e != nil {
		panic(e)
	}
	accessToken, accessTokenSecret, e := oauth1Config.AccessToken(session.Values["request_token"].(string), session.Values["request_token_secret"].(string), r.URL.Query().Get("oauth_verifier"))
	if e != nil {
		panic(e)
	}
	http.Redirect(w, r,
		session.Values["redirect_uri"].(string)+
			"?access_token="+accessToken+","+accessTokenSecret+
			"&state="+session.Values["state"].(string)+
			"&client_id="+session.Values["client_id"].(string)+
			"&response_type=Bearer",
		http.StatusFound)
}
