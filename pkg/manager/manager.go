// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"context"
	"github.com/onosproject/ran-simulator/pkg/mobility"
	"github.com/onosproject/ran-simulator/pkg/store/routes"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	cellapi "github.com/onosproject/ran-simulator/pkg/api/cells"
	metricsapi "github.com/onosproject/ran-simulator/pkg/api/metrics"
	modelapi "github.com/onosproject/ran-simulator/pkg/api/model"
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
	ModelName           string
	MetricName          string
	HOLogic             string
}

// NewManager creates a new manager
func NewManager(config *Config) (*Manager, error) {
	log.Info("Creating Manager")

	modelPluginRegistry := modelplugins.NewModelRegistry()
	for _, smp := range config.ServiceModelPlugins {
		if _, _, err := modelPluginRegistry.RegisterModelPlugin(smp); err != nil {
			log.Error(err)
		}
	}

	mgr := &Manager{
		config:              *config,
		agents:              nil,
		model:               &model.Model{},
		modelPluginRegistry: modelPluginRegistry,
	}

	return mgr, nil
}

// Manager is a manager for the E2T service
type Manager struct {
	modelapi.ManagementDelegate
	config              Config
	agents              *agents.E2Agents
	model               *model.Model
	modelPluginRegistry modelplugins.ModelRegistry
	server              *northbound.Server
	nodeStore           nodes.Store
	cellStore           cells.Store
	ueStore             ues.Store
	routeStore          routes.Store
	metricsStore        metrics.Store
	mobilityDriver      mobility.Driver
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
	err := model.Load(m.model, m.config.ModelName)
	if err != nil {
		log.Error(err)
		return err
	}

	m.initModelStores()
	m.initMetricStore()

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

	m.mobilityDriver = mobility.NewMobilityDriver(m.cellStore, m.routeStore, m.ueStore, m.model.APIKey, m.config.HOLogic)
	// TODO: Make initial speeds configurable
	m.mobilityDriver.GenerateRoutes(context.Background(), 30000, 160000, 20000)
	m.mobilityDriver.Start(context.Background())

	return nil
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
	m.stopE2Agents()
	m.stopNorthboundServer()
	m.mobilityDriver.Stop()
}

func (m *Manager) initModelStores() {
	// Create the node registry primed with the pre-loaded nodes
	m.nodeStore = nodes.NewNodeRegistry(m.model.Nodes)

	// Create the cell registry primed with the pre-loaded cells
	m.cellStore = cells.NewCellRegistry(m.model.Cells, m.nodeStore)

	// Create the UE registry primed with the specified number of UEs
	m.ueStore = ues.NewUERegistry(m.model.UECount, m.cellStore)

	// Create an empty route registry
	m.routeStore = routes.NewRouteRegistry()
}

func (m *Manager) initMetricStore() {
	// Create store for tracking arbitrary metrics and attributes for nodes, cells and UEs
	m.metricsStore = metrics.NewMetricsStore()

	// Load additional initial use-case data; ignore errors
	_ = pciload.LoadPCIMetrics(m.metricsStore, m.config.MetricName)
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
	m.server.AddService(modelapi.NewService(m))

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

func (m *Manager) stopE2Agents() {
	_ = m.agents.Stop()
}

func (m *Manager) stopNorthboundServer() {
	m.server.Stop()
}

// PauseAndClear pauses simulation and clears the model
func (m *Manager) PauseAndClear(ctx context.Context) {
	log.Info("Pausing RAN simulator...")
	m.stopE2Agents()
	m.nodeStore.Clear(ctx)
	m.cellStore.Clear(ctx)
	m.metricsStore.Clear(ctx)
}

// LoadModel loads the new model into the simulator
func (m *Manager) LoadModel(ctx context.Context, data []byte) error {
	m.model = &model.Model{}
	if err := model.LoadConfigFromBytes(m.model, data); err != nil {
		return err
	}
	m.initModelStores()
	return nil
}

// LoadMetrics loads new metrics into the simulator
func (m *Manager) LoadMetrics(ctx context.Context, name string, data []byte) error {
	if name == "rc.pci" {
		if err := pciload.LoadPCIMetricsData(m.metricsStore, data); err != nil {
			return err
		}
	}
	return nil
}

// Resume resume the simulation
func (m *Manager) Resume(ctx context.Context) {
	log.Info("Resuming RAN simulator...")
	go func() {
		time.Sleep(1 * time.Second)
		log.Info("Restarting NBI...")
		m.stopNorthboundServer()
		_ = m.startNorthboundServer()
	}()
	_ = m.startE2Agents()
}
