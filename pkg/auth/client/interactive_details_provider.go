package client

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

// InteractiveDetailsProvider prompts the user for the OAuth client ID and secret
type InteractiveDetailsProvider struct {
}

// NewInteractiveDetailsProvider returns a new InteractiveDetailsProvider
func NewInteractiveDetailsProvider() *InteractiveDetailsProvider {
	return &InteractiveDetailsProvider{}
}

func (p *InteractiveDetailsProvider) GetDetails(ctx context.Context) (*Details, error) {
	clientID, err := p.getClientID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client ID: %w", err)
	}

	clientSecret, err := p.getClientSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get client secret: %w", err)
	}

	return &Details{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

func (*InteractiveDetailsProvider) getClientID() (string, error) {
	clientIDPrompt := &promptui.Prompt{
		Label: "OAuth Client ID",
	}

	clientID, err := clientIDPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get client ID: %w", err)
	}

	return clientID, nil
}

func (*InteractiveDetailsProvider) getClientSecret() (string, error) {
	clientSecretPrompt := &promptui.Prompt{
		Label: "OAuth Client Secret",
	}

	clientSecret, err := clientSecretPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get client secret: %w", err)
	}

	return clientSecret, nil
}
