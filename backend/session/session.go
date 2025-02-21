package session

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/teamhanko/hanko/config"
	hankoJwk "github.com/teamhanko/hanko/crypto/jwk"
	hankoJwt "github.com/teamhanko/hanko/crypto/jwt"
	"net/http"
	"time"
)

type Manager interface {
	GenerateJWT(uuid.UUID) (string, error)
	Verify(string) (jwt.Token, error)
	GenerateCookie(userId uuid.UUID) (*http.Cookie, error)
}

// Manager is used to create and verify session JWTs
type manager struct {
	jwtGenerator  hankoJwt.Generator
	sessionLength time.Duration
	cookieConfig  cookieConfig
}

type cookieConfig struct {
	Domain   string
	HttpOnly bool
	SameSite http.SameSite
}

// NewManager returns a new Manager which will be used to create and verify sessions JWTs
func NewManager(jwkManager hankoJwk.Manager, config config.Session) (Manager, error) {
	signatureKey, err := jwkManager.GetSigningKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create session generator: %w", err)
	}
	verificationKeys, err := jwkManager.GetPublicKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to create session generator: %w", err)
	}
	g, err := hankoJwt.NewGenerator(signatureKey, verificationKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to create session generator: %w", err)
	}

	duration, _ := time.ParseDuration(config.Lifespan) // error can be ignored, value is checked in config validation
	sameSite := http.SameSite(0)
	switch config.Cookie.SameSite {
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}
	return &manager{
		jwtGenerator:  g,
		sessionLength: duration,
		cookieConfig: cookieConfig{
			Domain:   config.Cookie.Domain,
			HttpOnly: config.Cookie.HttpOnly,
			SameSite: sameSite,
		},
	}, nil
}

// GenerateJWT creates a new session JWT for the given user
func (g *manager) GenerateJWT(userId uuid.UUID) (string, error) {
	issuedAt := time.Now()
	expiration := issuedAt.Add(g.sessionLength)

	token := jwt.New()
	_ = token.Set(jwt.SubjectKey, userId.String())
	_ = token.Set(jwt.IssuedAtKey, issuedAt)
	_ = token.Set(jwt.ExpirationKey, expiration)
	//_ = token.Set(jwt.AudienceKey, []string{"http://localhost"})

	signed, err := g.jwtGenerator.Sign(token)
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

// Verify verifies the given JWT and returns a parsed one if verification was successful
func (g *manager) Verify(token string) (jwt.Token, error) {
	parsedToken, err := g.jwtGenerator.Verify([]byte(token))
	if err != nil {
		return nil, fmt.Errorf("failed to verify session token: %w", err)
	}

	return parsedToken, nil
}

// GenerateCookie creates a new session cookie for the given user
func (g *manager) GenerateCookie(userId uuid.UUID) (*http.Cookie, error) {
	jwt, err := g.GenerateJWT(userId)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     "hanko",
		Value:    jwt,
		Domain:   g.cookieConfig.Domain,
		Path:     "/",
		Secure:   true,
		HttpOnly: g.cookieConfig.HttpOnly,
		SameSite: g.cookieConfig.SameSite,
	}, nil
}
