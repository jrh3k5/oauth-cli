package client

import "context"

// DetailsProvider provides the OAuth client ID and secret
type DetailsProvider interface {
	// GetDetails returns the OAuth client ID and secret
	GetDetails(ctx context.Context) (*Details, error)
}
