package traefik_safeline

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	t1k "github.com/chaitin/t1k-go"
)

// Package example a example plugin.

// Config the plugin configuration.
type Config struct {
	// Addr is the address for the detector
	Addr     string `yaml:"addr"`
	PoolSize int    `yaml:"pool_size"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Addr:     "",
		PoolSize: 100,
	}
}

// Safeline a plugin.
type Safeline struct {
	next   http.Handler
	server *t1k.Server
	name   string
	config *Config
	logger *log.Logger
	mu     sync.Mutex
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := log.New(os.Stdout, "safeline", log.LstdFlags)
	logger.Printf("config: %+v", config)
	return &Safeline{
		next:   next,
		name:   name,
		config: config,
		logger: logger,
	}, nil
}

func (s *Safeline) initServer() error {
	if s.server != nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server == nil {
		server, err := t1k.NewWithPoolSize(s.config.Addr, s.config.PoolSize)
		if err != nil {
			return err
		}
		s.server = server
	}
	return nil
}

func (s *Safeline) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Printf("panic: %s", r)
		}
	}()
	if err := s.initServer(); err != nil {
		s.logger.Printf("error in initServer: %s", err)
		s.next.ServeHTTP(rw, req)
		return
	}
	rw.Header().Set("X-Chaitin-waf", "safeline")
	result, err := s.server.DetectHttpRequest(req)
	if err != nil {
		s.logger.Printf("error in detection: \n%+v\n", err)
		s.next.ServeHTTP(rw, req)
		return
	}
	if result.Blocked() {
		rw.WriteHeader(result.StatusCode())
		msg := fmt.Sprintf(`{"code": %d, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "%s"}`,
			result.StatusCode(),
			result.EventID(),
		)
		_, _ = rw.Write([]byte(msg))
		return
	}
	s.next.ServeHTTP(rw, req)
	//rw.WriteHeader(http.StatusForbidden)
	//_, _ = rw.Write([]byte("Inject by safeline\n"))
}
