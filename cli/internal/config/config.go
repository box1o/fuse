package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultAPIURL = "https://mback.teckstate.com"
	EnvAPIURL     = "FUSE_API_URL"
)

type Config struct {
	APIURL string
}

func Load() (Config, error) {
	apiURL := strings.TrimSpace(os.Getenv(EnvAPIURL))
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	if err := validateAPIURL(apiURL); err != nil {
		return Config{}, err
	}

	return Config{APIURL: strings.TrimRight(apiURL, "/")}, nil
}

func validateAPIURL(apiURL string) error {
	parsed, err := url.ParseRequestURI(apiURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("invalid Fuse API URL %q", apiURL)
	}

	return nil
}
