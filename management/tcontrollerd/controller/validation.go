package controller

import (
	"fmt"
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
)

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
