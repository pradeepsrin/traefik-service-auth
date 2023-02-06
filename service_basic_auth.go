// Package plugindemo a demo plugin.
package plugindemo

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	//Headers         map[string]string   `json:"headers,omitempty"`
	ClientHeaderKey string              `json:"client_header_key,omitempty"`
	Secrets         map[string][]string `json:"secrets,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Secrets:         make(map[string][]string),
		ClientHeaderKey: "",
	}
}

type ServiceBasicAuth struct {
	next     http.Handler
	config   *Config
	name     string
	template *template.Template
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.ClientHeaderKey == "" {
		return nil, fmt.Errorf("client_header_key cannot be empty/omited")
	}

	return &ServiceBasicAuth{
		config:   config,
		next:     next,
		name:     name,
		template: template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *ServiceBasicAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	givenClientId := req.Header.Get(a.config.ClientHeaderKey)
	if givenClientId == "" {
		errorOut(rw, "Client Header is not set")
		return
	}

	givenSecret := req.Header.Get("Authorization")
	if givenSecret == "" {
		errorOut(rw, "Not Authorized")
		return
	}

	clientIdSecrets := a.config.Secrets[givenClientId]
	if clientIdSecrets == nil {
		errorOut(rw, "Not Authorized")
		return
	}

	for _, secret := range clientIdSecrets {
		if givenSecret == secret {
			a.next.ServeHTTP(rw, req)
			return
		}
	}

	errorOut(rw, "Not Authorized")
	return
}

func errorOut(rw http.ResponseWriter, s string) {
	http.Error(rw, s, http.StatusForbidden)
}
