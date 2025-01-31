package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jrh3k5/oauth-cli/pkg/auth"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

const (
	oauthClientID     = "000000"
	oauthClientSecret = "999999"
)

func main() {
	ctx := context.Background()

	runOAuthServer()

	// Wait for the server to be available
	for {
		_, err := http.Get("http://localhost:8080/authorize")
		if err == nil {
			break
		}
	}

	log.Println("OAuth server started")

	oauthToken, err := auth.DefaultGetOAuthToken(ctx,
		"http://localhost:8080/authorize",
		"http://localhost:8080/token",
	)

	if err != nil {
		panic(fmt.Errorf("failed to get OAuth token: %w", err))
	}

	log.Printf("OAuth token: %s\n", oauthToken.AccessToken)
}

func runOAuthServer() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set(oauthClientID, &models.Client{
		ID:     oauthClientID,
		Secret: oauthClientSecret,
		Domain: "http://localhost:54520",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.UserAuthorizationHandler = func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return oauthClientID, nil
	}

	srv.SetInternalErrorHandler(func(err error) *errors.Response {
		return &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	go http.ListenAndServe("0.0.0.0:8080", nil)
}
