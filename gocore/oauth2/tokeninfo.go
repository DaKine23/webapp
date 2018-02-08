package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var cache = make(map[string]TokenInfo)

// TokenInfo contains all information you can get from token info service
type TokenInfo struct {
	AccessToken string   `json:"access_token"`
	ClientID    string   `json:"client_id"`
	ExpiresIn   int      `json:"expires_in"`
	GrantType   string   `json:"grant_type"`
	Realm       string   `json:"realm"`
	Scope       []string `json:"scope"`
	TokenType   string   `json:"token_type"`
	UID         string   `json:"uid"`
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

var clnt Client = http.DefaultClient

//Init is used to define a costum client  *(if not called http.DefaultClient is used)*
func Init(cl Client) {
	clnt = cl
}

// RetrieveTokenInfo gets TokenInfo from token info service (https)
// **tokenInfoService**  is *mandatory* and should contain the host like : info.services.auth.zalando.com ;
// **token** is *mandatory* and should contain the token ;
// use Init(cl Client) to define a costum client  *(if not http.DefaultClient is used)*
func RetrieveTokenInfo(tokenInfoService, token string) (*TokenInfo, error) {

	ti, ok := cache[token]
	if ok {
		return &ti, nil
	}

	//check mandatory parameter
	switch {
	case len(tokenInfoService) == 0:
		return nil, errors.New("tokenInfoService parameter was empty")
	case len(token) == 0:
		return nil, errors.New("token parameter was empty")
	}

	//create URL
	tisu := url.URL{
		Scheme: "https",
		Host:   tokenInfoService,
		Path:   "oauth2/tokeninfo",
	}
	//add encoded query
	query := tisu.Query()
	query.Set("access_token", token)
	tisu.RawQuery = query.Encode()

	//call
	req, _ := http.NewRequest("GET", tisu.String(), nil) //err only occures on empty url so its ignored here

	return retrieveScopes(clnt.Do(req))

}

func retrieveScopes(res *http.Response, err error) (*TokenInfo, error) {

	//check for initial errors
	if err != nil {
		return nil, err
	}

	//check status code for != 2xx
	if (res.StatusCode / 100) != 2 {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("Request for scopes returned a not 2XX response. status code : %v, status : %s, header : %v, body : %s", res.StatusCode, res.Status, res.Header, string(body))
	}

	//decode body
	var tokeninfo TokenInfo
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&tokeninfo); err != nil {
		return nil, err
	}

	//cache tokens in memory
	cache[tokeninfo.AccessToken] = tokeninfo
	go func() {
		timer := time.NewTimer(time.Second * time.Duration(tokeninfo.ExpiresIn))
		<-timer.C
		delete(cache, tokeninfo.AccessToken)
	}()

	return &tokeninfo, nil

}
