package detection

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type ResultObjective int

const (
	RO_REQUEST  ResultObjective = 0
	RO_RESPONSE ResultObjective = 1
)

type Result struct {
	Objective   ResultObjective
	Head        byte
	Body        []byte
	Alog        []byte
	ExtraHeader []byte
	ExtraBody   []byte
	T1KContext  []byte
	Cookie      []byte
	WebLog      []byte
}

func (r *Result) Passed() bool {
	return r.Head == '.'
}

func (r *Result) Blocked() bool {
	return !r.Passed()
}

func (r *Result) StatusCode() int {
	str := string(r.Body)
	if str == "" {
		return http.StatusForbidden
	}
	code, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("t1k convert status code failed: %v", err)
		return http.StatusForbidden
	}
	return code
}

func (r *Result) BlockMessage() map[string]interface{} {
	return map[string]interface{}{
		"status":   r.StatusCode(),
		"success":  false,
		"message":  "blocked by Chaitin SafeLine Web Application Firewall",
		"event_id": r.EventID(),
	}
}

func (r *Result) EventID() string {
	extra := string(r.ExtraBody)
	if extra == "" {
		return ""
	}
	// <!-- event_id: e1impksyjq0gl92le6odi0fnobi270cj -->
	re, err := regexp.Compile(`<\!--\s*event_id:\s*([a-zA-Z0-9]+)\s*-->\s*`)
	if err != nil {
		log.Printf("t1k compile regexp failed: %v", err)
		return ""
	}
	matches := re.FindStringSubmatch(extra)
	if len(matches) < 2 {
		log.Printf("t1k regexp not match event id: %s", extra)
		return ""
	}
	return matches[1]
}
