package client

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	var oauthClientSecret string

	for _, osArg := range os.Args {
		switch {
		case strings.HasPrefix(osArg, "--oauth-client-id="):
			splitArg := strings.Split(osArg, "=")
			if len(splitArg) != 2 {
				return "", "", fmt.Errorf("invalid flag: %s", osArg)
			}

			oauthClientID = splitArg[1]
		case strings.HasPrefix(osArg, "--oauth-client-secret="):
			splitArg := strings.Split(osArg, "=")
			if len(splitArg) != 2 {
				return "", "", fmt.Errorf("invalid flag: %s", osArg)
			}

			oauthClientSecret = splitArg[1]
		}
	}

	if oauthClientID == "" {
		return "", "", fmt.Errorf("--oauth-client-id is required")
	}

	if oauthClientSecret == "" {
		return "", "", fmt.Errorf("--oauth-client-secret is required")
	}

	return oauthClientID, oauthClientSecret, nil
}
