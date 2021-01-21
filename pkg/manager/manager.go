// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
	"github.com/onosproject/ran-simulator/pkg/model"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath   string
	KeyPath  string
	CertPath string
	GRPCPort int
}

// NewManager creates a new manager
func NewManager(config *Config) (*Manager, error) {
	log.Info("Creating Manager")
	mgr := &Manager{
		config: *config,
		agents: nil,
		model:  &model.Model{},
	}

	return mgr, nil
}

// Manager is a manager for the E2T service
type Manager struct {
	config Config
	agents *e2agent.E2Agents
	model  *model.Model
	server *northbound.Server
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager:", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	// Create the E2 agents for all simulated nodes and specified controllers
	err := model.Load(m.model)
	if err != nil {
		log.Error(err)
		return err
	}

	m.agents = e2agent.NewE2Agents(m.model)
	// Start the E2 agents
	err = m.agents.Start()
	if err != nil {
		return err
	}

	// Start gRPC server
	return m.startNorthboundServer()
}

// startSouthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	m.server = northbound.NewServer(northbound.NewServerCfg(
		m.config.CAPath,
		m.config.KeyPath,
		m.config.CertPath,
		int16(m.config.GRPCPort),
		true,
		northbound.SecurityConfig{}))
	m.server.AddService(logging.Service{})

	doneCh := make(chan error)
	go func() {
		err := m.server.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
	_ = m.agents.Stop()
	m.stopNorthboundServer()
}

func (m *Manager) stopNorthboundServer() {
	// TODO implementation requires ability to actually stop the server
}
