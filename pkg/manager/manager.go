// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	cellapi "github.com/onosproject/ran-simulator/pkg/api/cells"
	metricsapi "github.com/onosproject/ran-simulator/pkg/api/metrics"
	nodeapi "github.com/onosproject/ran-simulator/pkg/api/nodes"
	"github.com/onosproject/ran-simulator/pkg/api/trafficsim"
	"github.com/onosproject/ran-simulator/pkg/e2agent/agents"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/rc/pciload"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/metrics"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath              string
	KeyPath             string
	CertPath            string
	GRPCPort            int
	ServiceModelPlugins []string
}

// NewManager creates a new manager
func NewManager(config *Config) (*Manager, error) {
	log.Info("Creating Manager")

	modelPluginRegistry := modelplugins.ModelPluginRegistry{
		ModelPlugins: make(map[modelplugins.ModelFullName]modelplugins.ModelPlugin),
	}
	for _, smp := range config.ServiceModelPlugins {
		if _, _, err := modelPluginRegistry.RegisterModelPlugin(smp); err != nil {
			log.Error(err)
		}
	}

	mgr := &Manager{
		config:              *config,
		agents:              nil,
		model:               &model.Model{},
		modelPluginRegistry: &modelPluginRegistry,
	}

	return mgr, nil
}

// Manager is a manager for the E2T service
type Manager struct {
	config              Config
	agents              *agents.E2Agents
	model               *model.Model
	modelPluginRegistry *modelplugins.ModelPluginRegistry
	server              *northbound.Server
	nodeStore           nodes.Store
	cellStore           cells.Store
	ueStore             ues.Store
	metricsStore        metrics.Store
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Error("Unable to run Manager:", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	// Load the model data
	err := model.Load(m.model)
	if err != nil {
		log.Error(err)
		return err
	}

	// Create the node registry primed with the pre-loaded nodes
	m.nodeStore = nodes.NewNodeRegistry(m.model.Nodes)

	// Create the cell registry primed with the pre-loaded cells
	m.cellStore = cells.NewCellRegistry(m.model.Cells, m.nodeStore)

	// Create the UE registry primed with the specified number of UEs
	m.ueStore = ues.NewUERegistry(m.model.UECount, m.cellStore)

	// Create store for tracking arbitrary metrics and attributes for nodes, cells and UEs
	m.metricsStore = metrics.NewMetricsStore()

	// Load additional initial use-case data; ignore errors
	_ = pciload.LoadPCIMetrics(m.metricsStore)

	// Start gRPC server
	err = m.startNorthboundServer()
	if err != nil {
		return err
	}
	// Start E2 agents
	err = m.startE2Agents()
	if err != nil {
		return err
	}

	return nil
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
	m.server.AddService(nodeapi.NewService(m.nodeStore, m.model.PlmnID))
	m.server.AddService(cellapi.NewService(m.cellStore))
	m.server.AddService(trafficsim.NewService(m.model, m.cellStore, m.ueStore))
	m.server.AddService(metricsapi.NewService(m.metricsStore))

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

func (m *Manager) startE2Agents() error {
	// Create the E2 agents for all simulated nodes and specified controllers
	var err error
	m.agents, err = agents.NewE2Agents(m.model, m.modelPluginRegistry,
		m.nodeStore, m.ueStore, m.cellStore, m.metricsStore)
	if err != nil {
		log.Error(err)
		return err
	}
	// Start the E2 agents
	err = m.agents.Start()
	if err != nil {
		return err
	}

	return nil
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
