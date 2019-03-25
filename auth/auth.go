package auth

import (
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"os"
)

type Authenticator struct {
	firebaseApp *firebase.App
}

func NewAuthenticator() (*Authenticator, error) {
	b := []byte(os.Getenv("FIREBASE_ACCOUNT"))
	opt := option.WithCredentialsJSON(b)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return &Authenticator{app}, nil
}

func (auth *Authenticator) Authenticate(ctx context.Context, id string) (*string, error) {
	client, err := auth.firebaseApp.Auth(ctx)
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, id)
	if err != nil {
		return nil, err
	}
	return &token.UID, nil
}

