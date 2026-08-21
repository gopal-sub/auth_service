package google

import (
	"net/url"
	"os"
)

type GoogleProvider struct {
	
	cliendID string
	cliendSecret string
	redirectURI string
	scope string
}


func NewGoogleProvider() *GoogleProvider{
	return &GoogleProvider{
		cliendID: os.Getenv("GOOGLE_CLIENT_ID"),
		cliendSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		redirectURI: "http://localhost:3000/callback",
		scope: "openid",
	}
}

func (g *GoogleProvider) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", g.cliendID)
	params.Add("redirect_uri", g.redirectURI)
	params.Add("scope", g.scope)
	params.Add("response_type", "code")
	authorisationEndPoint := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
	return  authorisationEndPoint

	
}

// func (g *GoogleProvider) HandleCallback (code string) (string, error) {
// 	queryParams := .URL.Query()

// 	code := queryParams.Get("code")
// 	err := queryParams.Get("error")
// 	if err != "" {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"error": "shit went down",
// 		})
// 	}
// 	tokenEndpoint := "https://oauth2.googleapis.com/token"
// 	redirectURI := "http://localhost:3000/callback"

// 	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
// 	clientID := os.Getenv("GOOGLE_CLIENT_ID")

// 	data := url.Values{}

// 	data.Add("client_id", clientID)
// 	data.Add("client_secret", clientSecret)
// 	data.Add("code", code)
// 	data.Add("grant_type", "authorization_code")
// 	data.Add("redirect_uri", redirectURI)

// 	res, _ := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

// }