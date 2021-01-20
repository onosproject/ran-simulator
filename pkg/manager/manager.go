// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/ran-simulator/pkg/e2agent"
	"github.com/onosproject/ran-simulator/pkg/model"
	smregistry "github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"gopkg.in/yaml.v2"
	"os"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath   string `yaml:"caPath"`
	KeyPath  string `yaml:"keyPath"`
	CertPath string `yaml:"certPath"`
	GRPCPort int    `yaml:"grpcPort"`

	// Path to the YAML file describing the environment model
	ModelPath string `yaml:"modelPath"`
}

// NewManager creates a new manager
func NewManager(config *Config) (*Manager, error) {
	log.Info("Creating Manager")
	mgr := &Manager{
		config:   *config,
		model:    model.NewModel(),
		agents:   nil,
		registry: smregistry.NewServiceModelRegistry(),
	}

	err := mgr.model.Load(config.ModelPath)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

// LoadConfig from the specified YAML file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	config := &Config{}
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// Manager is a manager for the E2T service
type Manager struct {
	config   Config
	model    *model.Model
	agents   *e2agent.E2Agents
	registry *smregistry.ServiceModelRegistry
	server   *northbound.Server
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
	// Create the E2 agents for all simulated nodes and specified controllers
	m.agents = e2agent.NewE2Agents(m.model.Nodes.GetAll(), m.registry, m.model.Controllers)

	// S tart the E2 agents
	err := m.agents.Start()
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
