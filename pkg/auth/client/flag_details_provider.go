package client

import (
	"context"
	"flag"
	"fmt"
)

// FlagDetailsProvider provides the OAuth client ID and secret by reading arguments provided as flags.
// This reads --oauth-client-id and --oauth-client-secret.
type FlagDetailsProvider struct {
}

// NewFlagDetailsProvider returns a new FlagDetailsProvider.
func NewFlagDetailsProvider() *FlagDetailsProvider {
	return &FlagDetailsProvider{}
}

func (p *FlagDetailsProvider) GetDetails(ctx context.Context) (*Details, error) {
	clientID, clientSecret, err := getFlagValues()
	if err != nil {
		return nil, fmt.Errorf("failed to get flag values: %w", err)
	}

	return &Details{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

func getFlagValues() (string, string, error) {
	var oauthClientID string
	flag.StringVar(&oauthClientID, "oauth-client-id", "", "OAuth client ID")

	var oauthClientSecret string
	flag.StringVar(&oauthClientSecret, "oauth-client-secret", "", "OAuth client secret")

	flag.Parse()

	if oauthClientID == "" {
		return "", "", fmt.Errorf("--oauth-client-id is required")
	}

	if oauthClientSecret == "" {
		return "", "", fmt.Errorf("--oauth-client-secret is required")
	}

	return oauthClientID, oauthClientSecret, nil
}
