package google

import (
	"os"
	"gopal-sub/auth_service/internal/auth/oauth"
)



func NewGoogleProvider() *oauth.Provider{
	return &oauth.Provider{
		CliendID: os.Getenv("GOOGLE_CLIENT_ID"),
		CliendSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURI: "http://localhost:3000/google/callback",
		Scope: "openid",
		AuthorisationEndPointBase: "https://accounts.google.com/o/oauth2/v2/auth?",
		TokenEndpoint: "https://oauth2.googleapis.com/token",
	}
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