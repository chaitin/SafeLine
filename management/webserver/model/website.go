package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/datatypes"
)

const (
	// maxPort is the largest TCP port that may be listened on.
	maxPort = 65535

	// maxServerNameLength and maxLabelLength follow the DNS limits.
	maxServerNameLength = 253
	maxLabelLength      = 63
)

var (
	// portPattern matches the decimal notation of a TCP port.
	portPattern = regexp.MustCompile(`^[0-9]{1,5}$`)

	// serverNameLabelPattern matches one label of a literal host name. The
	// accepted character set is deliberately narrow: the value ends up in an
	// nginx "server_name" directive, so anything that could terminate the
	// directive (or start another one) must not pass.
	serverNameLabelPattern = regexp.MustCompile(`^[A-Za-z0-9_]([A-Za-z0-9_-]*[A-Za-z0-9_])?$`)
)

type Website struct {
	Base
	Comment     string         `gorm:"comment"           json:"comment"`
	ServerNames datatypes.JSON `gorm:"server_names"      json:"server_names"`
	Ports       datatypes.JSON `gorm:"ports"             json:"ports"`
	Upstreams   datatypes.JSON `gorm:"upstreams"         json:"upstreams"`

	CertFilename string `gorm:"cert_filename"    json:"cert_filename"`
	KeyFilename  string `gorm:"key_filename"     json:"key_filename"`

	IsEnabled bool `gorm:"is_enabled;default=true"       json:"is_enabled"`
}

// Validate checks the fields that tcontrollerd renders into the nginx
// configuration of the site.
//
// Ports and server names are interpolated into nginx directives, so they have
// to be plain ports and host names: a value carrying nginx syntax must be
// rejected here, before it is stored and pushed to the engine.
func (w *Website) Validate() error {
	serverNames, err := unmarshalStringList(w.ServerNames)
	if err != nil {
		return fmt.Errorf("invalid server_names: %v", err)
	}

	ports, err := unmarshalStringList(w.Ports)
	if err != nil {
		return fmt.Errorf("invalid ports: %v", err)
	}

	for _, serverName := range serverNames {
		if err = ValidateServerName(serverName); err != nil {
			return err
		}
	}

	for _, port := range ports {
		if err = ValidatePort(port); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePort accepts the decimal notation of a TCP port.
func ValidatePort(port string) error {
	if !portPattern.MatchString(port) {
		return fmt.Errorf("invalid port %q: only digits are allowed", port)
	}

	value, err := strconv.Atoi(port)
	if err != nil || value < 1 || value > maxPort {
		return fmt.Errorf("invalid port %q: must be between 1 and %d", port, maxPort)
	}

	return nil
}

// ValidateServerName accepts a literal host name, optionally with a wildcard
// as its first label, or "*" for the default site.
func ValidateServerName(name string) error {
	// "*" (and the empty string, which the console sends for the same purpose)
	// makes the site the default server and is rendered as "_".
	if name == "*" || name == "" {
		return nil
	}

	if len(name) > maxServerNameLength {
		return fmt.Errorf("invalid server name %q: longer than %d characters", name, maxServerNameLength)
	}

	for i, label := range strings.Split(name, ".") {
		if label == "*" {
			// nginx only accepts a wildcard as the left most label.
			if i != 0 {
				return fmt.Errorf("invalid server name %q: wildcard must be the first label", name)
			}
			continue
		}

		if len(label) > maxLabelLength || !serverNameLabelPattern.MatchString(label) {
			return fmt.Errorf("invalid server name %q: %q is not a valid host name label", name, label)
		}
	}

	return nil
}

// unmarshalStringList decodes one of the JSON list columns of a website.
func unmarshalStringList(raw datatypes.JSON) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}

	return values, nil
}
