package auth
import "context"

type AuthService interface {
    InitiateRegistrationFlow(ctx context.Context) (map[string]interface{}, error)
    CompleteRegistration(ctx context.Context, flowID string, req RegistrationRequest) (*Identity, error)
    ValidateSession(ctx context.Context, sessionToken string) (*Identity, error)
    CheckCollegeAccess(identity *Identity, collegeID string) bool
    HasRole(identity *Identity, role string) bool
    GetPublicURL()string 
}