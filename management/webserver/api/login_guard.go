package api

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// loginFailureThreshold is the number of failed attempts that are tolerated
	// before the caller has to wait.
	loginFailureThreshold = 5

	// loginFailureWindow is how long a failed attempt is remembered.
	loginFailureWindow = 10 * time.Minute

	// loginMaxLockout caps the wait that is imposed on a caller that keeps
	// failing. A six digit passcode stays out of reach of a brute force attack
	// with this delay, even over a long period.
	loginMaxLockout = 5 * time.Minute

	// loginMaxEntries bounds the memory the throttle may use. The keys carry a
	// client address, and a client can choose as many of those as it likes.
	loginMaxEntries = 1024
)

// loginAttempt is the state of one throttling key.
type loginAttempt struct {
	failures    int
	lastFailure time.Time
}

// loginGuard throttles the login endpoint per account and per client address.
//
// The management console has a single account and the passcode is its only
// credential, so a caller that is allowed to try passcodes without limit can
// simply enumerate them.
type loginGuard struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt

	// now is a seam for the tests.
	now func() time.Time
}

func newLoginGuard() *loginGuard {
	return &loginGuard{
		attempts: make(map[string]*loginAttempt),
		now:      time.Now,
	}
}

// accountKeyPrefix marks the throttle key that every caller of an account
// shares, as opposed to the keys that belong to a single client address.
const accountKeyPrefix = "account:"

// isAccountKey reports whether a throttle key protects an account rather than a
// single client address.
func isAccountKey(key string) bool {
	return strings.HasPrefix(key, accountKeyPrefix)
}

// loginKeys returns the throttle keys of a login attempt.
func loginKeys(clientIP, username string) []string {
	keys := []string{accountKeyPrefix + username}
	if clientIP != "" {
		keys = append(keys, fmt.Sprintf("ip:%s", clientIP))
	}

	return keys
}

// lockout reports how long the caller has to wait, and the key that caused it.
func (g *loginGuard) lockout(keys []string) (string, time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.now()
	g.prune(now)

	var (
		blockedKey string
		wait       time.Duration
	)

	for _, key := range keys {
		attempt, ok := g.attempts[key]
		if !ok || attempt.failures < loginFailureThreshold {
			continue
		}

		if remaining := loginBackoff(attempt.failures) - now.Sub(attempt.lastFailure); remaining > wait {
			blockedKey, wait = key, remaining
		}
	}

	return blockedKey, wait
}

// fail records an authentication failure.
func (g *loginGuard) fail(keys []string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := g.now()
	g.prune(now)

	for _, key := range keys {
		attempt, ok := g.attempts[key]
		if !ok || now.Sub(attempt.lastFailure) > loginFailureWindow {
			attempt = &loginAttempt{}
			g.attempts[key] = attempt
		}

		attempt.failures++
		attempt.lastFailure = now
	}
}

// reset forgets the failures of a key, which is what a successful login does.
func (g *loginGuard) reset(keys []string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, key := range keys {
		delete(g.attempts, key)
	}
}

// prune drops records that are no longer useful and keeps the map bounded.
func (g *loginGuard) prune(now time.Time) {
	for key, attempt := range g.attempts {
		if now.Sub(attempt.lastFailure) > loginFailureWindow+loginMaxLockout {
			delete(g.attempts, key)
		}
	}

	// The map has to stay bounded, because a client can invent client addresses
	// faster than the failures expire. Dropping an address key only weakens the
	// throttle of a caller that stopped trying, but an account key is never
	// dropped: it is the one that keeps the passcode out of reach, and being
	// able to reset it by filling the map with invented addresses would make
	// the whole throttle useless.
	for key := range g.attempts {
		if len(g.attempts) <= loginMaxEntries {
			break
		}
		if isAccountKey(key) {
			continue
		}
		delete(g.attempts, key)
	}
}

// loginBackoff returns how long a caller that failed n times in a row has to
// wait after its last failure. The delay doubles per failure and is capped.
func loginBackoff(failures int) time.Duration {
	shift := failures - loginFailureThreshold
	if shift < 0 {
		shift = 0
	}
	if shift > 8 {
		shift = 8
	}

	backoff := time.Second << uint(shift)
	if backoff > loginMaxLockout {
		backoff = loginMaxLockout
	}

	return backoff
}
