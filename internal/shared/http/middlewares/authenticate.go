package middlewares

import (
	"context"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/mariocoski/address-service/internal/config"
	"github.com/sirupsen/logrus"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	// Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func Authenticate(config config.Config) func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse(config.Auth0Domain)
	if err != nil {
		logrus.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Auth0Audience},
		// validator.WithCustomClaims(
		// 	func() validator.CustomClaims {
		// 		return &CustomClaims{}
		// 	},
		// ),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		logrus.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		logrus.Infof("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthenticated"}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}
