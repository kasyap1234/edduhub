package auth

import (
	"context"
	"eduhub/server/internal/services/auth"
)

type AuthService interface {
    InitiateRegistrationFlow(ctx context.Context) (map[string]interface{}, error)
    CompleteRegistration(ctx context.Context, flowID string, req RegistrationRequest) (*Identity, error)
    ValidateSession(ctx context.Context, sessionToken string) (*Identity, error)
    CheckCollegeAccess(identity *Identity, collegeID string) bool
    HasRole(identity *Identity, role string) bool
    HasPermission(ctx context.Context,identity *Identity,action,resource string)(bool,error)
    AssignRole(ctx context.Context,identityID string,role string)error 
    RemoveRole(ctx context.Context,identityID string, role string)error 
    AddPermission(ctx context.Context, identityID, action, resource string) error
    RemovePermission(ctx context.Context, identityID, action, resource string) error
    
    GetPublicURL()string 

}


type authService struct {
    Auth *kratosService
    AuthZ  *ketoService
}

func NewAuthService(kratos *kratosService,keto *ketoService)*authService{
    return &authService{
       Auth: kratos , 
       AuthZ: keto,
       
    }
}

func(a*authService)InitiateRegistrationFlow(ctx context.Context)(map[string]interface{},error){
   return  a.Auth.InitiateRegistrationFlow(ctx)

}
