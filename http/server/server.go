package server

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HttpServer struct{ http http.Server }

func New(opts ...Option) *HttpServer {
	s := &HttpServer{
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

type Option func(*HttpServer)

func WithMaxHeaderBytes(i uint) Option {
	return func(s *HttpServer) { s.http.MaxHeaderBytes = int(i) }
}

func WithReadTimeout(t time.Duration) Option {
	return func(s *HttpServer) { s.http.ReadTimeout = t }
}

func WithReadHeaderTimeout(t time.Duration) Option {
	return func(s *HttpServer) { s.http.ReadHeaderTimeout = t }
}

func WithWriteTimeout(t time.Duration) Option {
	return func(s *HttpServer) { s.http.WriteTimeout = t }
}

func WithIdleTimeout(t time.Duration) Option {
	return func(s *HttpServer) { s.http.IdleTimeout = t }
}

type TLSVersion int

const (
	TlsOld TLSVersion = iota
	TlsIntermediate
	TlsModern
)

func WithTLSConfig(i TLSVersion) Option {
	return func(s *HttpServer) {
		switch i {
		case TlsOld:
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
		case TlsIntermediate:
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
		case TlsModern:
			s.http.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS13,
			}
		default:
			panic("invalid argument")
		}
	}
}
