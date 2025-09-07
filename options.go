package gap

import (
	"net/url"
)

type Options struct {
	Backend    *url.URL `short:"b" required:"" env:"GAP_BACKEND" help:"Backend URL."`
	Port       uint     `short:"p" required:"" env:"GAP_PORT" help:"Listening port."`
	HeaderName string   `short:"n" required:"" env:"GAP_HEADER" help:"Header name to pass the access token."`
	AllowList  []string `short:"e" required:"" env:"GAP_ALLOW_LIST" help:"Allowed email list that may contain wildcards."`
}
