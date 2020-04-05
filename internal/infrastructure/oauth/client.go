package oauth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const clientDefaultTimeout = time.Second * 10

func NewClient(ctx context.Context, token string) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := oauth2.NewClient(ctx, ts)

	client.Timeout = clientDefaultTimeout

	return client
}
