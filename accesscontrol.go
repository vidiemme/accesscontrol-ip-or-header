// Package accesscontrol provides an access control plugin for Traefik.
package accesscontrol

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Config is the plugin configuration.
type Config struct {
	Whitelist []string `json:"whitelist,omitempty"`
	HeaderKey string   `json:"headerKey,omitempty"`
	HeaderValue string `json:"headerValue,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Whitelist: []string{},
		HeaderKey: "",
		HeaderValue: "",
	}
}

// AccessControl is the plugin structure.
type AccessControl struct {
	next       http.Handler
	whitelist  map[string]struct{}
	headerKey  string
	headerValue string
}

// New creates a new AccessControl plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Whitelist) == 0 && (config.HeaderKey == "" || config.HeaderValue == "") {
		return nil, fmt.Errorf("either whitelist or header validation must be configured")
	}

	whitelist := make(map[string]struct{}, len(config.Whitelist))
	for _, ip := range config.Whitelist {
		whitelist[ip] = struct{}{}
	}

	return &AccessControl{
		next:       next,
		whitelist:  whitelist,
		headerKey:  config.HeaderKey,
		headerValue: config.HeaderValue,
	}, nil
}

func (a *AccessControl) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	clientIP := strings.Split(req.RemoteAddr, ":")[0]

	// Check if the client's IP is in the whitelist.
	if _, allowed := a.whitelist[clientIP]; allowed {
		a.next.ServeHTTP(rw, req)
		return
	}

	// Check if the required header is present and matches the expected value.
	if req.Header.Get(a.headerKey) == a.headerValue {
		a.next.ServeHTTP(rw, req)
		return
	}

	// Deny access if neither condition is satisfied.
	http.Error(rw, "Access denied", http.StatusForbidden)
}
