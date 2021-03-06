package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// CompleteAuthentication exchanges the given code for token
func CompleteAuthentication(w http.ResponseWriter, r *http.Request, p *models.PluginContext) {

	authenticatedUserID := r.Header.Get("Mattermost-User-ID")
	if authenticatedUserID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	conf := p.OauthConfig

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")

	storedState, e := p.API.KVGet(state)
	if e != nil {
		mlog.Error(e.Error())
		http.Error(w, "missing stored state", http.StatusBadRequest)
		return
	} else if string(storedState) != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	userID := strings.Split(state, "_")[1]

	p.API.KVDelete(state)

	if userID != authenticatedUserID {
		http.Error(w, "Not authorized, incorrect user", http.StatusUnauthorized)
		return
	}

	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", p.UserInfoEndpoint, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tok.AccessToken))
	res, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1<<20))

	user := &models.UserInfo{}
	err = json.Unmarshal(body, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenBytes, err := json.Marshal(tok)
	p.API.KVSet(models.NowTokenKeyPrefix+userID, tokenBytes)

	p.API.KVSet(models.NowUserIDPrefix+user.Result.ID, []byte(userID))

	html := `
		<!DOCTYPE html>
		<html>
			<head>
				<script>
					window.close();
				</script>
			</head>
			<body>
				<p>Completed connecting to ServiceNow. Please close this window.</p>
			</body>
		</html>
		`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Authorize redirects to the authorize endpoint
func Authorize(w http.ResponseWriter, r *http.Request, p *models.PluginContext) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	oauthConfig := p.OauthConfig

	state := fmt.Sprintf("%v_%v", model.NewId()[0:15], userID)

	p.API.KVSet(state, []byte(state))

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	mlog.Info("Authorize Url: " + url)

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
