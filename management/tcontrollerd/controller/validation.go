package controller

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

	// certFilenamePattern matches the file name of a certificate or a key. The
	// name is rendered into "ssl_certificate" and "ssl_certificate_key", and
	// it is joined to the certificate directory, so only a plain file name
	// without a path separator, a control character or any nginx syntax may
	// pass.
	certFilenamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

	// upstreamSchemePattern matches the two schemes an upstream may use. The
	// scheme is rendered into "proxy_pass", so anything else is rejected.
	upstreamSchemePattern = regexp.MustCompile(`^https?$`)

	// upstreamHostPattern matches the host of an upstream: a host name, an IPv4
	// address or a bracketed IPv6 address, with an optional port. The host is
	// rendered into the "server" directive of the site, so a value carrying
	// nginx syntax (a space, a semicolon, a "$", a quote) must not pass.
	upstreamHostPattern = regexp.MustCompile(`^(?:\[[0-9A-Fa-f:.]+\]|[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?\.?)(?::([0-9]{1,5}))?$`)
)

// maxCertFilenameLength bounds a certificate file name.
const maxCertFilenameLength = 255

// validateCertFilename accepts the name of a file in the certificate directory.
//
// The management server validates the same value before it stores it, but this
// side cannot rely on that: the configuration is pushed over the network, and a
// name carrying nginx syntax would turn into directives that run commands.
func validateCertFilename(name string) error {
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

// validatePort accepts the decimal notation of a TCP port. It rejects any
// other value because the port is interpolated into a "listen" directive.
func validatePort(port string) error {
	if !portPattern.MatchString(port) {
		return fmt.Errorf("invalid port %q: only digits are allowed", port)
	}

	value, err := strconv.Atoi(port)
	if err != nil || value < 1 || value > maxPort {
		return fmt.Errorf("invalid port %q: must be between 1 and %d", port, maxPort)
	}

	return nil
}

// validateUpstream accepts the address of a site upstream.
//
// The management server validates the same value before it stores it, but this
// side cannot rely on that: the site configuration is pushed over the network,
// and the scheme and the host of an upstream are rendered into the "proxy_pass"
// and "server" directives of the site.
func validateUpstream(upstream string) error {
	parsed, err := url.Parse(upstream)
	if err != nil {
		return fmt.Errorf("invalid upstream %q: %v", upstream, err)
	}

	if !upstreamSchemePattern.MatchString(parsed.Scheme) {
		return fmt.Errorf("invalid upstream %q: the scheme has to be http or https", upstream)
	}

	host := parsed.Host
	if parsed.Scheme == HttpsScheme && parsed.Port() == "" {
		// generateNginxConfig appends the default HTTPS port to this value
		// before it renders it, so validate the value it renders.
		host = host + ":" + DefaultHttpsPort
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

// validateServerName accepts a literal host name, optionally with a wildcard
// as its first label. It rejects any other value because the name is
// interpolated into a "server_name" directive.
func validateServerName(name string) error {
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
