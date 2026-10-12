package oauth

import (
	"net/url"
)
type Provider struct {
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