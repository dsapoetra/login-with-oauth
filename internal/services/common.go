package services

import (
	"login-with-oauth/internal/helpers/pages"
	"login-with-oauth/internal/logger"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)

	if err != nil {
		logger.Log.Error("Error parsing OAuth URL: " + err.Error())
	}

	logger.Log.Info("URL: " + URL.String())
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)

	URL.RawQuery = parameters.Encode()

	url := URL.String()
	logger.Log.Info(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pages.IndexPage))
}
