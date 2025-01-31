package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/int128/oauth2cli"
	"github.com/jrh3k5/oauth-cli/pkg/auth/client"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

// DefaultOAuthServerPort is the default port for the OAuth server to listen for successful authentication
const DefaultOAuthServerPort = 54520

// Logger is the type of logger used during OAuth token retrieval
type Logger func(format string, args ...any)

// TokenOptions contains options for GetOAuthToken
type TokenOptions struct {
	logger     Logger
	serverPort int
	scopes     []string
}

// TokenOption is an option for GetOAuthToken
type TokenOption func(*TokenOptions)

// WithLogger sets the logger
func WithLogger(logger Logger) TokenOption {
	return func(opts *TokenOptions) {
		opts.logger = logger
	}
}

// WithOAuthServerPort sets the port on which the OAuth server should run
func WithOAuthServerPort(port int) TokenOption {
	return func(opts *TokenOptions) {
		opts.serverPort = port
	}
}

// WithScopes sets the OAuth scopes to be requested for the token.
func WithScopes(scopes ...string) TokenOption {
	return func(opts *TokenOptions) {
		opts.scopes = append(opts.scopes, scopes...)
	}
}

// DefaultGetOAuthToken provides a default behavior for calling GetOAuthToken, using the --interative flag to determine whether to use interactive mode
// or non-interactive mode.
func DefaultGetOAuthToken(ctx context.Context, authURL string, tokenURL string, opts ...TokenOption) (*oauth2.Token, error) {
	isInteractive := false
	for _, osArg := range os.Args {
		isInteractive = isInteractive || osArg == "--interactive"
		if isInteractive {
			break
		}
	}

	var clientDetailsProvider client.DetailsProvider
	if isInteractive {
		clientDetailsProvider = client.NewInteractiveDetailsProvider()
	} else {
		clientDetailsProvider = client.NewFlagDetailsProvider()
	}

	return GetOAuthToken(ctx, authURL, tokenURL, clientDetailsProvider, opts...)
}

// GetOAuthToken gets the OAuth token from the given OAuth server using the given client details provider.
func GetOAuthToken(ctx context.Context, authURL string, tokenURL string, clientDetailsProvider client.DetailsProvider, opts ...TokenOption) (*oauth2.Token, error) {
	tokenOptions := &TokenOptions{
		serverPort: DefaultOAuthServerPort,
		logger: func(string, ...any) {
			// deliberately no-op
		},
	}

	for _, opt := range opts {
		opt(tokenOptions)
	}

	localServerURLChan := make(chan string)
	defer close(localServerURLChan)

	clientDetails, err := clientDetailsProvider.GetDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get client details: %w", err)
	}

	cliConfig := oauth2cli.Config{
		OAuth2Config: oauth2.Config{
			ClientID:     clientDetails.ClientID,
			ClientSecret: clientDetails.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
			Scopes:      tokenOptions.scopes,
			RedirectURL: fmt.Sprintf("http://127.0.0.1:%d", tokenOptions.serverPort),
		},
		LocalServerBindAddress: []string{fmt.Sprintf("127.0.0.1:%d", tokenOptions.serverPort)},
		LocalServerReadyChan:   localServerURLChan,
		Logf:                   tokenOptions.logger,
	}

	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		select {
		case url := <-localServerURLChan:
			if err := browser.OpenURL(url); err != nil {
				tokenOptions.logger("could not open the browser: %s", err)
			}
			return nil
		case <-errGroupCtx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})

	tokenReturn := make(chan *oauth2.Token, 1)
	errGroup.Go(func() error {
		token, err := oauth2cli.GetToken(errGroupCtx, cliConfig)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}

		tokenReturn <- token

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	close(tokenReturn)

	return <-tokenReturn, nil
}
