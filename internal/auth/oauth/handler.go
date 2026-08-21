package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"io"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	svc Oauthsvc
}


func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.svc.GetAuthLoginURL("google")
	fmt.Println(url)
}
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {


}


// -------------------

func (h *Handler) OauthLogin(w http.ResponseWriter, r *http.Request){

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	// 
	redirectURI := "http://localhost:3000/callback"


	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "openid profile email")
	params.Add("response_type", "code")

	authorisationEndPoint := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()


	http.Redirect(w, r, authorisationEndPoint, http.StatusMovedPermanently)


}

func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request){
	
	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	err := queryParams.Get("error")
	if err != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "shit went down",
		})
	}
	tokenEndpoint := "https://oauth2.googleapis.com/token"
	redirectURI := "http://localhost:3000/callback"

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	data := url.Values{}

	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", redirectURI)

	res, _ := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	body, _ := io.ReadAll(res.Body)
	type id_token struct{
		AccessToken string `json:"access_token"`  
		ExpiresIn int `json:"expires_in"` 
		Score string `json:"scope"` 
		TokenType string `json:"token_type"` 
		IdToken string `json:"id_token"` 
	}

	id_token_get := id_token{}

	json.Unmarshal(body, &id_token_get)

	claims := jwt.MapClaims{}

	token, _, _ := jwt.NewParser().ParseUnverified(id_token_get.IdToken, claims)
	
	fmt.Println(claims["email"])
	fmt.Println(token)



}