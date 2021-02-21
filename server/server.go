package server

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Server struct{ http http.Server }

func NewHTTPServer(opts ...ServerOption) *Server {
	s := &Server{
		http: http.Server{
			TLSConfig:   &tls.Config{MinVersion: tls.VersionTLS13},
			ReadTimeout: time.Second * 5,
			// ReadHeaderTimeout: time.Second * 0,
			// WriteTimeout:      time.Second * 0,
			IdleTimeout:    time.Second * 90,
			MaxHeaderBytes: 2 << 20, // 2MB
		},
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

type ServerOption func(*Server)

func WithMaxHeaderBytes(i uint) ServerOption {
	return func(s *Server) { s.http.MaxHeaderBytes = int(i) }
}

func WithReadTimeout(t time.Duration) ServerOption {
	return func(s *Server) { s.http.ReadTimeout = t }
}

func WithReadHeaderTimeout(t time.Duration) ServerOption {
	return func(s *Server) { s.http.ReadHeaderTimeout = t }
}

func WithWriteTimeout(t time.Duration) ServerOption {
	return func(s *Server) { s.http.WriteTimeout = t }
}

func WithIdleTimeout(t time.Duration) ServerOption {
	return func(s *Server) { s.http.IdleTimeout = t }
}

type TLSVersion int

const (
	TlsConfigOld TLSVersion = iota
	TlsConfigIntermediate
	TlsConfigModern
)

func WithTLSConfig(i TLSVersion) ServerOption {
	return func(s *Server) {
		switch i {
		case TlsConfigOld:
			s.http.TLSConfig = &tls.Config{
				MinVersion:               tls.VersionTLS10,
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			}
		case TlsConfigIntermediate:
			s.http.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			}
		case TlsConfigModern:
			s.http.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS13,
			}
		}
	}
}
