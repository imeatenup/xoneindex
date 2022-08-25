package intercept

import (
	"context"
	"net/http"

	"github.com/xnuc/xoneindex/config"
	"github.com/xnuc/xoneindex/log"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var Token *oauth2.Token

type OauthClient struct{}

func (i *OauthClient) PreHandle(w http.ResponseWriter, r *http.Request) bool {
	var err error
	defer func() {
		if err != nil {
			log.Errorf(r.Context(), "OauthClient.PreHandle err{%+v}", err)
		}
	}()
	log.Debugf(r.Context(), "OauthClient.PreHandle")
	if Token != nil { // valid token continue
		return true
	}

	r.ParseForm() // callback to new token
	state := r.Form.Get("state")
	if state != config.STATE {
		http.Error(w, "authorized by state", http.StatusInternalServerError)
		return false
	}
	if r.URL.Path == config.CALLBACK && len(r.Form.Get("code")) != 0 {
		Token, err = i.callback(r.Context(), r.Form.Get("code"))
		if err != nil {
			http.Error(w, "authorized by code", http.StatusInternalServerError)
			return false
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return false
	}

	redirect := (&oauth2.Config{ // invalid token to authorize
		ClientID:    config.CLIENT_ID,
		Scopes:      config.Scopes,
		RedirectURL: config.REDIRECT_URI,
		Endpoint:    microsoft.AzureADEndpoint("consumers"),
	}).AuthCodeURL(state)
	http.Redirect(w, r, redirect, http.StatusFound)
	return false
}

func (i *OauthClient) PostHandle(_ http.ResponseWriter, r *http.Request) bool {
	log.Debugf(r.Context(), "OauthClient.PostHandle")
	return true
}

func (i *OauthClient) callback(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := (&oauth2.Config{
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
		Scopes:       config.Scopes,
		RedirectURL:  config.REDIRECT_URI,
		Endpoint:     microsoft.AzureADEndpoint("consumers"),
	}).Exchange(ctx, code)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return token, nil
}
