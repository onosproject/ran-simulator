// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package simmanager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/agent"
	smregistry "github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"net"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath     string
	KeyPath    string
	CertPath   string
	GRPCPort   int
	E2THost    string
	E2SCTPPort int
}

// NewManager creates a new manager
func NewManager(config Config) *Manager {
	log.Info("Creating Manager")
	return &Manager{
		Config: config,
	}
}

// Manager is a manager for the E2T service
type Manager struct {
	Config  Config
	e2agent agent.Agent
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	registry := smregistry.NewServiceModelRegistry()
	ips, err := net.LookupIP(m.Config.E2THost)
	if err != nil {
		return err
	}
	addr := ips[0].String()
	m.e2agent = agent.NewE2Agent(registry, addr, m.Config.E2SCTPPort)
	agentErr := m.e2agent.Start()
	if agentErr != nil {
		return agentErr
	}

	nbErr := m.startNorthboundServer()
	if nbErr != nil {
		return nbErr
	}
	return nbErr
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
	_ = m.e2agent.Stop()
}

// startSouthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.Config.CAPath,
		m.Config.KeyPath,
		m.Config.CertPath,
		int16(m.Config.GRPCPort),
		true,
		northbound.SecurityConfig{}))
	s.AddService(logging.Service{})

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}
