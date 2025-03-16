package auth

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

const (
    PUBLIC_URL = "KRATOS_PUBLIC_URL"
    ADMIN_URL  = "KRATOS_ADMIN_URL"
)




type KratosService struct {
    PublicURL  string
    AdminURL   string
    HTTPClient *http.Client
}

type Identity struct {
    ID     string `json:"id"`
    Traits struct {
        Email    string `json:"email"`
        Name     struct {
            First string `json:"first"`
            Last  string `json:"last"`
        } `json:"name"`
        College struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        } `json:"college"`
        Role string `json:"role"`
    } `json:"traits"`
}

type RegistrationRequest struct {
    Method   string `json:"method"`
    Password string `json:"password"`
    Traits   struct {
        Email    string `json:"email"`
        Name     struct {
            First string `json:"first"`
            Last  string `json:"last"`
        } `json:"name"`
        College struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        } `json:"college"`
        Role string `json:"role"`
    } `json:"traits"`
}

func NewKratosService() *KratosService {
    return &KratosService{
        PublicURL:  os.Getenv(PUBLIC_URL),
        AdminURL:   os.Getenv(ADMIN_URL),
        HTTPClient: &http.Client{},
    }
}

// InitiateRegistrationFlow starts registration process
func (k *KratosService) InitiateRegistrationFlow(ctx context.Context) (map[string]interface{}, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", 
        fmt.Sprintf("%s/self-service/registration/api", k.PublicURL), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create registration request: %w", err)
    }

    resp, err := k.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute registration request: %w", err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode registration response: %w", err)
    }

    return result, nil
}

// CompleteRegistration submits registration data
func (k *KratosService) CompleteRegistration(ctx context.Context, flowID string, req RegistrationRequest) (*Identity, error) {
    data, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal registration data: %w", err)
    }

    url := fmt.Sprintf("%s/self-service/registration?flow=%s", k.PublicURL, flowID)
    request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
    if err != nil {
        return nil, fmt.Errorf("failed to create registration completion request: %w", err)
    }

    request.Header.Set("Content-Type", "application/json")
    resp, err := k.HTTPClient.Do(request)
    if err != nil {
        return nil, fmt.Errorf("failed to complete registration: %w", err)
    }
    defer resp.Body.Close()

    var identity Identity
    if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
        return nil, fmt.Errorf("failed to decode identity: %w", err)
    }

    return &identity, nil
}

// ValidateSession checks if session is valid and returns identity
func (k *KratosService) ValidateSession(ctx context.Context, sessionToken string) (*Identity, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", 
        fmt.Sprintf("%s/sessions/whoami", k.PublicURL), nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("X-Session-Token", sessionToken)
    resp, err := k.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("invalid session")
    }

    var result struct {
        Identity Identity `json:"identity"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result.Identity, nil
}

// CheckCollegeAccess verifies if user belongs to specific college
func (k *KratosService) CheckCollegeAccess(identity *Identity, collegeID string) bool {
    return identity.Traits.College.ID == collegeID
}

// HasRole checks if user has specific role
func (k *KratosService) HasRole(identity *Identity, role string) bool {
    return identity.Traits.Role == role
}