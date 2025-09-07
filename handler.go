package gap

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/winebarrel/sglob"
	oauth2opt "google.golang.org/api/option"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

type AuthHandler struct {
	Options *Options
	Oauth2  *Oauth2Client
	Proxy   func(http.ResponseWriter, *http.Request)
}

func NewAuthHandler(options *Options, oauth2opts ...oauth2opt.ClientOption) (*AuthHandler, error) {
	oc, err := NewOauth2Client(oauth2opts...)

	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(options.Backend)

	handler := &AuthHandler{
		Options: options,
		Oauth2:  oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			r.Host = options.Backend.Host
			proxy.ServeHTTP(w, r)
		},
	}

	return handler, nil
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(h.Options.HeaderName)

	if token == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	ti, err := h.Oauth2.Tokeninfo(token)

	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if ti.Email == "" {
		http.Error(w, "cannot get email", http.StatusForbidden)
		return
	}

	for _, allowed := range h.Options.AllowList {
		if sglob.Match(allowed, ti.Email) {
			log.Printf("%s -> %s", ti.Email, r.URL.Path)
			h.Proxy(w, r)
			return
		}
	}

	http.Error(w, "not allowed", http.StatusForbidden)
}
