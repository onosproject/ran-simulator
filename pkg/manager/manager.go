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
func NewManager(config *Config) *Manager {
	log.Info("Creating Manager")
	return &Manager{
		Config:   *config,
		Model:    model.NewModel(),
		Agents:   nil,
		Registry: smregistry.NewServiceModelRegistry(),
	}
}

// LoadConfig from the specified YAML file
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// Manager is a manager for the E2T service
type Manager struct {
	Config   Config
	Model    *model.Model
	Agents   *e2agent.E2Agents
	Registry *smregistry.ServiceModelRegistry
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
	// Start E2 agents watchdog
	err := m.startE2Agents()
	if err != nil {
		return err
	}

	// Start gRPC server
	return m.startNorthboundServer()
}

func (m *Manager) startE2Agents() error {
	//ips, err := net.LookupIP(m.Config.E2THost)
	//if err != nil {
	//	return err
	//}
	//addr := ips[0].String()
	//m.e2agent = e2agent.NewE2Agent(registry, addr, m.Config.E2SCTPPort)
	//agentErr := m.e2agent.Start()
	//if agentErr != nil {
	//	return agentErr
	//}
	return nil
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

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
	stopE2Agents()
	stopNorthboundServer()
}

func stopE2Agents() {
	//_ = m.e2agent.Stop()
}

func stopNorthboundServer() {

}
