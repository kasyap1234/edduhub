package config

type AuthConfig struct {
	Domain string 
	ClientID string 
	ClientSecret string 
	Scopes string 
	// more fields 
}

func (a *AuthConfig)LoadAuthConfig(*AuthConfig){
	
	// type AuthConfig , 
	

}
