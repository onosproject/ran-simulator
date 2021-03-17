// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package northbound houses implementations of various application-oriented interfaces
// for the ONOS configuration subsystem.
package northbound

import (
	"crypto/tls"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/grpcinterceptors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound")

// Service provides service-specific registration for grpc services.
type Service interface {
	Register(s *grpc.Server)
}

// Server provides NB gNMI server for onos-config.
type Server struct {
	cfg      *ServerConfig
	services []Service
	server   *grpc.Server
}

// SecurityConfig security configuration
type SecurityConfig struct {
	AuthenticationEnabled bool
	AuthorizationEnabled  bool
}

// ServerConfig comprises a set of server configuration options.
type ServerConfig struct {
	CaPath      *string
	KeyPath     *string
	CertPath    *string
	Port        int16
	Insecure    bool
	SecurityCfg *SecurityConfig
}

// NewServer initializes gNMI server using the supplied configuration.
func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		services: []Service{},
		cfg:      cfg,
	}
}

// NewServerConfig creates a server config created with the specified end-point security details.
// Deprecated: Use NewServerCfg instead
func NewServerConfig(caPath string, keyPath string, certPath string, port int16, secure bool) *ServerConfig {
	return &ServerConfig{
		Port:        port,
		Insecure:    secure,
		CaPath:      &caPath,
		KeyPath:     &keyPath,
		CertPath:    &certPath,
		SecurityCfg: &SecurityConfig{},
	}
}

// NewServerCfg creates a server config created with the specified end-point security details.
func NewServerCfg(caPath string, keyPath string, certPath string, port int16, secure bool, secCfg SecurityConfig) *ServerConfig {
	return &ServerConfig{
		Port:        port,
		Insecure:    secure,
		CaPath:      &caPath,
		KeyPath:     &keyPath,
		CertPath:    &certPath,
		SecurityCfg: &secCfg,
	}
}

// AddService adds a Service to the server to be registered on Serve.
func (s *Server) AddService(r Service) {
	s.services = append(s.services, r)
}

// Serve starts the NB gNMI server.
func (s *Server) Serve(started func(string)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}
	tlsCfg := &tls.Config{}

	if *s.cfg.CertPath == "" && *s.cfg.KeyPath == "" {
		// Load default Certificates
		clientCerts, err := tls.X509KeyPair([]byte(certs.DefaultLocalhostCrt), []byte(certs.DefaultLocalhostKey))
		if err != nil {
			log.Error("Error loading default certs")
			return err
		}
		tlsCfg.Certificates = []tls.Certificate{clientCerts}
	} else {
		log.Infof("Loading certs: %s %s", *s.cfg.CertPath, *s.cfg.KeyPath)
		clientCerts, err := tls.LoadX509KeyPair(*s.cfg.CertPath, *s.cfg.KeyPath)
		if err != nil {
			log.Info("Error loading default certs")
		}
		tlsCfg.Certificates = []tls.Certificate{clientCerts}
	}

	if s.cfg.Insecure {
		// RequestClientCert will ask client for a certificate but won't
		// require it to proceed. If certificate is provided, it will be
		// verified.
		tlsCfg.ClientAuth = tls.RequestClientCert
	} else {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if *s.cfg.CaPath == "" {
		log.Info("Loading default CA onfca")
		tlsCfg.ClientCAs, err = certs.GetCertPoolDefault()
	} else {
		tlsCfg.ClientCAs, err = certs.GetCertPool(*s.cfg.CaPath)
	}
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsCfg))}
	if s.cfg.SecurityCfg.AuthenticationEnabled {
		log.Info("Authentication Enabled")
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(grpcinterceptors.AuthenticationInterceptor),
			)))
		opts = append(opts, grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_auth.StreamServerInterceptor(grpcinterceptors.AuthenticationInterceptor))))

	}

	s.server = grpc.NewServer(opts...)
	for i := range s.services {
		s.services[i].Register(s.server)
	}
	started(lis.Addr().String())

	log.Infof("Starting RPC server on address: %s", lis.Addr().String())
	return s.server.Serve(lis)
}

// Stop stops the server.
func (s *Server) Stop() {
	s.server.Stop()
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}
