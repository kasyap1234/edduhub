
package config 

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
  )
  

type SuperTokensConfig struct {
ConnectionURI string 
APIKey string 
AppName string 
APIDomain string 
WebDomain string 
// apiBasePath string 
// websiteBasePath string 
}




  func InitSuperTokens(config SuperTokensConfig) {
	  apiBasePath := "/auth"
	  websiteBasePath := "/auth"
	  err := supertokens.Init(supertokens.TypeInput{
		  Supertokens: &supertokens.ConnectionInfo{
		 
		  ConnectionURI: config.ConnectionURI,
		  APIKey: config.APIKey,
		  },
		  AppInfo: supertokens.AppInfo{
			AppName: config.AppName,
			APIDomain: config.APIDomain,
			WebsiteDomain: config.WebDomain,
				  APIBasePath: &apiBasePath,
				  WebsiteBasePath: &websiteBasePath,
		  },
		  RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		  },
	  })
  
	  if err != nil {
		panic(err.Error())
	  }
  }

