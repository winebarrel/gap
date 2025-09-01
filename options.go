package gap

import (
	"net/url"
)

type Options struct {
	Backend *url.URL `short:"b" required:"" env:"GAP_BACKEND" help:"Backend URL to proxy."`
	Port    uint     `short:"p" required:"" env:"GAP_PORT" help:"Listening port."`
	Domain  string   `short:"d" required:"" env:"GAP_DOMAIN" help:"Allowed email domain."`
}
