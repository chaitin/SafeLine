package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
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

	// defaultHttpsPort is the port tcontrollerd appends to an upstream that
	// names https without a port.
	defaultHttpsPort = "443"
)

var (
	// portPattern matches the decimal notation of a TCP port.
	portPattern = regexp.MustCompile(`^[0-9]{1,5}$`)

	// serverNameLabelPattern matches one label of a literal host name. The
	// accepted character set is deliberately narrow: the value ends up in an
	// nginx "server_name" directive, so anything that could terminate the
	// directive (or start another one) must not pass.
	serverNameLabelPattern = regexp.MustCompile(`^[A-Za-z0-9_]([A-Za-z0-9_-]*[A-Za-z0-9_])?$`)

	// certFilenamePattern matches the file name of an uploaded certificate or
	// key. The name is interpolated into "ssl_certificate" and
	// "ssl_certificate_key", and it is joined to the certificate directory, so
	// only a plain file name without a path separator, a control character or
	// any nginx syntax may pass.
	certFilenamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

	// upstreamSchemePattern matches the two schemes an upstream may use. The
	// scheme is interpolated into "proxy_pass", so anything else is rejected.
	upstreamSchemePattern = regexp.MustCompile(`^https?$`)

	// upstreamHostPattern matches the host of an upstream: a host name, an IPv4
	// address or a bracketed IPv6 address, with an optional port. The host is
	// interpolated into the "server" directive of the site, so a value carrying
	// nginx syntax (a space, a semicolon, a "$", a quote) must not pass.
	upstreamHostPattern = regexp.MustCompile(`^(?:\[[0-9A-Fa-f:.]+\]|[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?\.?)(?::([0-9]{1,5}))?$`)
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

	upstreams, err := unmarshalStringList(w.Upstreams)
	if err != nil {
		return fmt.Errorf("invalid upstreams: %v", err)
	}

	for _, upstream := range upstreams {
		if err = ValidateUpstream(upstream); err != nil {
			return err
		}
	}

	if err = ValidateCertFilename(w.CertFilename); err != nil {
		return err
	}

	if err = ValidateCertFilename(w.KeyFilename); err != nil {
		return err
	}

	return nil
}

// ValidateUpstream accepts the address of a site upstream.
//
// The scheme and the host of this value are interpolated into the "proxy_pass"
// and "server" directives of the site, so a value carrying nginx syntax must be
// rejected here, before it is stored and pushed to the engine.
func ValidateUpstream(upstream string) error {
	parsed, err := url.Parse(upstream)
	if err != nil {
		return fmt.Errorf("invalid upstream %q: %v", upstream, err)
	}

	if !upstreamSchemePattern.MatchString(parsed.Scheme) {
		return fmt.Errorf("invalid upstream %q: the scheme has to be http or https", upstream)
	}

	host := parsed.Host
	if parsed.Scheme == "https" && parsed.Port() == "" {
		// tcontrollerd appends the default HTTPS port to this value before it
		// renders it, so validate the value it renders.
		host = host + ":" + defaultHttpsPort
	}

	matches := upstreamHostPattern.FindStringSubmatch(host)
	if matches == nil {
		return fmt.Errorf("invalid upstream %q: %q is not a host name, an IP address or a host with a port", upstream, host)
	}

	if port := matches[1]; port != "" {
		value, err := strconv.Atoi(port)
		if err != nil || value < 1 || value > maxPort {
			return fmt.Errorf("invalid upstream %q: port has to be between 1 and %d", upstream, maxPort)
		}
	}

	return nil
}

// maxCertFilenameLength bounds a certificate file name.
const maxCertFilenameLength = 255

// ValidateCertFilename accepts the name of a file in the certificate directory.
//
// The console stores uploaded certificates under a generated name, so a value
// that is not a plain file name is either a mistake or an attempt to inject
// nginx directives (through the rendered ssl_certificate directives) or to
// address a file outside the certificate directory.
func ValidateCertFilename(name string) error {
	// An empty name means "this site is served over plain HTTP".
	if name == "" {
		return nil
	}

	if len(name) > maxCertFilenameLength {
		return fmt.Errorf("invalid certificate file name %q: longer than %d characters", name, maxCertFilenameLength)
	}

	if name != filepath.Base(name) || name == "." || name == ".." {
		return fmt.Errorf("invalid certificate file name %q: only a file name without a directory is allowed", name)
	}

	if !certFilenamePattern.MatchString(name) {
		return fmt.Errorf("invalid certificate file name %q: only letters, digits, dot, dash and underscore are allowed", name)
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
