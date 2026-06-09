package idp

type IDPAuthorizeURLResponse struct {
	AuthorizationURL string `json:"authorizationUrl"`
}

type IDPTokenExchangeRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type IDPSessionResponse struct {
	Message string `json:"message"`
}

type IDPLogoutResponse struct {
	Message string `json:"message"`
}

// IDPTokenResponse represents the response from IDP token endpoint
type IDPTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// IDPUserInfo represents user information from IDP userinfo endpoint
type IDPUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	SuffixName string `json:"name_suffix"`
}
