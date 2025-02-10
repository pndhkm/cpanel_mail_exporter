package utils

import (
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"golang.org/x/crypto/bcrypt"
)

type TLSConfig struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type WebConfig struct {
	TLSServerConfig TLSConfig         `yaml:"tls_server_config"`
	BasicAuthUsers  map[string]string `yaml:"basic_auth_users"`
}

// LoadWebConfig reads the YAML configuration file
func LoadWebConfig(path string) (*WebConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config WebConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// BasicAuth is a middleware for HTTP basic authentication
func BasicAuth(next http.Handler, users map[string]string) http.Handler {
	if len(users) == 0 {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !validateUser(username, password, users) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// validateUser checks if provided username and password match with stored bcrypt hash
func validateUser(username, password string, users map[string]string) bool {
	hashedPassword, exists := users[username]
	if !exists {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// IsTLSEnabled checks if both cert and key files are provided
func IsTLSEnabled(tlsConfig TLSConfig) bool {
	return tlsConfig.CertFile != "" && tlsConfig.KeyFile != ""
}
