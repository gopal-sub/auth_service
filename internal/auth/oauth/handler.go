package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"io"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	Provider Provider
	
}

func NewHandler(Provider Provider) *Handler{
	return &Handler{
		Provider: Provider,
	}
}

type Id_token struct{
	AccessToken string `json:"access_token"`  
	ExpiresIn int `json:"expires_in"` 
	Score string `json:"scope"` 
	TokenType string `json:"token_type"` 
	IdToken string `json:"id_token"` 
}


// -------------------

func (h *Handler) OauthLogin(w http.ResponseWriter, r *http.Request){

	clientID := h.Provider.CliendID
	redirectURI := h.Provider.RedirectURI


	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", h.Provider.Scope)
	params.Add("response_type", "code")

	authorisationEndPoint := h.Provider.AuthorisationEndPointBase + params.Encode()


	http.Redirect(w, r, authorisationEndPoint, http.StatusMovedPermanently)

}

func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request) string{
	//refactor this
	
	queryParams := r.URL.Query()

	code := queryParams.Get("code")[0]
	err := queryParams.Get("error")
	if err != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "oauth callbakc query param issue",
		})
	}
	return string(code)

}

func (h *Handler) GetToken(code string){

	data := url.Values{}

	data.Add("client_id", h.Provider.CliendID)
	data.Add("client_secret", h.Provider.CliendSecret)
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", h.Provider.RedirectURI)

	res, _ := http.Post(h.Provider.TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	body, _ := io.ReadAll(res.Body)
	

	IDToken := Id_token{}

	json.Unmarshal(body, &IDToken)

	claims := jwt.MapClaims{}
	fmt.Print(claims)

	token, _, _ := jwt.NewParser().ParseUnverified(IDToken.IdToken, claims)
	
	fmt.Println(claims["email"])
	fmt.Println(token)

}