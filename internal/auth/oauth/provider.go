package oauth

import (
	"net/url"
)

type ProviderEntity string

const (
	GOOGLE ProviderEntity = "GOOGLE"
	GITHUB ProviderEntity = "GITHUB"
)

var ValidProvider = map[ProviderEntity] bool {
	GOOGLE: true,
	GITHUB: false,
}

type Provider struct {
	ProviderEntity ProviderEntity
	CliendID string
	CliendSecret string
	RedirectURI string
	Scope string
	AuthorisationEndPointBase string
	TokenEndpoint string
}


func (p *Provider) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", p.CliendID)
	params.Add("redirect_uri", p.RedirectURI)
	params.Add("scope", p.Scope)
	params.Add("response_type", "code")
	authorisationEndPoint := p.AuthorisationEndPointBase + params.Encode()

	return  authorisationEndPoint
	
}