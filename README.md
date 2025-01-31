# OAuth CLI Tooling

This is utility used to capture an OAuth token for a CLI tooling. It supports an interactive mode (where the caller is asked to provide the OAuth client ID and secret) or non-interactive mode, where the client ID and secret can be supplied as arguments to the CLI tooling.

## Usage

An example of how to use this in your CLI tooling is:

```
package main

import (
    "context"
    "github.com/jrh3k5/oauth-cli/auth"
)

func main() {
    ctx := context.Background()

    oauthToken, err := auth.DefaultGetOAuthToken(ctx, "http://server.com/authorize", "http://server.com/token")
    if err != nil {
        // handle error
    }

    // use token
}
```

### Example Usage

This library provides a means of exercising the library. You can run it by one of two ways:

```
go run internal/examples/main.go --oauth-client-id=000000 --oauth-client-secret=999999
```

And, with `000000` as the client ID and `999999` as the client secret, this can be called interactively:

```
go run internal/examples/main.go --interactive
```

This example tool uses the OAuth Playground available [here](https://www.oauth.com/playground/). Register a client there and supply the client ID and secret as appropriate.