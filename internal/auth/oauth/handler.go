package oauth

import (
	"encoding/json"
	"fmt"
	"gopal-sub/auth_service/internal/redis"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
)

type Handler struct {
	Redis *redis.Repository
	
}

func NewHandler(Redis *redis.Repository) *Handler{
	return &Handler{
		Redis: Redis,
	}
}

// svc => provider ?? 

type Id_token struct{
	AccessToken string `json:"access_token"`  
	ExpiresIn int `json:"expires_in"` 
	Score string `json:"scope"` 
	TokenType string `json:"token_type"` 
	IdToken string `json:"id_token"` 
}




// -------------------

func (h *Handler) OauthCodeFlow(provider *Provider, w http.ResponseWriter, r *http.Request){

	clientID := provider.CliendID
	redirectURI := provider.RedirectURI
	state := randstr.String(10)

	redisKey := RedisKeyPrefix + state
	
	h.Redis.Insert(redisKey, string(provider.ProviderEntity))
	


	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", provider.Scope)
	params.Add("response_type", "code")
	params.Add("state", state)

	authorisationEndPoint := provider.AuthorisationEndPointBase + params.Encode()


	http.Redirect(w, r, authorisationEndPoint, http.StatusMovedPermanently)

}

//callback handler should be indpendent

func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request){
	//refactor this
	
	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	state := queryParams.Get("state")
	err := queryParams.Get("error")
	if err != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "oauth callback query param issue",
		})
	}

	if code != "" && state != "" {
		providerEntity, err := h.OauthStateValidator(state)

		token, err := GetToken(googleProvider, code, state)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "state invalid",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})

		
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "no code or state ",
	})
	
	
}

func (h *Handler) OauthStateValidator(state string) (ProviderEntity, error) {
	//check if the state is valid 
	redisKey := RedisKeyPrefix + state
	providerEntity ,err := h.Redis.Get(redisKey)
	if err != redis.KeydoesNotExist{
		return "",StateDoesNotExist 
	}
	if err != nil {
		return "", err
	}
	return ProviderEntity(providerEntity), nil

}

func GetToken(provider *Provider, code string, state string) (string, error){


	
	// if the control reaches here => state valid

	data := url.Values{}

	data.Add("client_id", provider.CliendID)
	data.Add("client_secret", provider.CliendSecret)
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", provider.RedirectURI)
	data.Add("state", state)

	res, _ := http.Post(provider.TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	body, _ := io.ReadAll(res.Body)
	

	IDToken := Id_token{}

	json.Unmarshal(body, &IDToken)

	claims := jwt.MapClaims{}
	fmt.Print(claims)

	token, _, _ := jwt.NewParser().ParseUnverified(IDToken.IdToken, claims)
	
	fmt.Println(claims["email"])
	fmt.Println(token)

	return "", nil

}