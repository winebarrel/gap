package gap

import (
	"context"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Oauth2Client struct {
	Service *oauth2.Service
}

func NewOauth2Client(opts ...option.ClientOption) (*Oauth2Client, error) {
	opts = append(opts, option.WithoutAuthentication())
	service, err := oauth2.NewService(context.Background(), opts...)

	if err != nil {
		return nil, err
	}

	client := &Oauth2Client{
		Service: service,
	}

	return client, nil
}

func (client *Oauth2Client) Tokeninfo(token string) (*oauth2.Tokeninfo, error) {
	ti, err := client.Service.Tokeninfo().AccessToken(token).Do()

	if err != nil {
		return nil, err
	}

	return ti, nil
}
