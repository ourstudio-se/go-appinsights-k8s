package appink8s

import (
	"sync"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
)

type initializer interface {
	ReadPropertySpec() (*runtimeSpec, error)
}

type kubernetesTelemetryClient struct {
	appinsights.TelemetryClient
	active      bool
	initializer initializer
	initialized bool
	lock        sync.RWMutex
	properties  map[string]string
}

func NewTelemetryClient(iKey string) appinsights.TelemetryClient {
	cfg := newK8sConfig()
	if !cfg.RunningInKubernetes() {
		return appinsights.NewTelemetryClient(iKey)
	}

	client, err := newK8sClient(cfg)
	if err != nil {
		return appinsights.NewTelemetryClient(iKey)
	}

	return &kubernetesTelemetryClient{
		TelemetryClient: appinsights.NewTelemetryClient(iKey),
		active:          true,
		initializer:     newK8sInitializer(client),
		initialized:     false,
		properties:      make(map[string]string),
	}
}

func (ktc *kubernetesTelemetryClient) apply(properties map[string]string) {
	if !ktc.initialized {
		ktc.initialize()
	}

	if !ktc.active {
		return
	}

	ktc.lock.RLock()
	defer ktc.lock.RUnlock()

	for k, v := range ktc.properties {
		properties[k] = v
	}
}

func (ktc *kubernetesTelemetryClient) initialize() {
	ktc.lock.Lock()
	defer ktc.lock.Unlock()

	spec, err := ktc.initializer.ReadPropertySpec()
	ktc.active = err == nil
	ktc.initialized = true
	ktc.properties = spec.ToPropertyMap()

	if err == nil && spec.DeploymentName != "" {
		ktc.Context().Tags.Cloud().SetRole(spec.DeploymentName)
		ktc.Context().Tags.Cloud().SetRoleInstance(spec.PodName)
	}
}

func (ktc *kubernetesTelemetryClient) Track(t appinsights.Telemetry) {
	ktc.apply(t.GetProperties())
	ktc.TelemetryClient.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackAvailability(name string, duration time.Duration, success bool) {
	t := appinsights.NewAvailabilityTelemetry(name, duration, success)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackEvent(name string) {
	t := appinsights.NewEventTelemetry(name)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackException(err interface{}) {
	t := appinsights.NewExceptionTelemetry(err)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackMetric(name string, value float64) {
	t := appinsights.NewMetricTelemetry(name, value)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackRemoteDependency(name, dependencyType, target string, success bool) {
	t := appinsights.NewRemoteDependencyTelemetry(name, dependencyType, target, success)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackRequest(method, uri string, duration time.Duration, responseCode string) {
	t := appinsights.NewRequestTelemetry(method, uri, duration, responseCode)
	ktc.Track(t)
}

func (ktc *kubernetesTelemetryClient) TrackTrace(name string, severity contracts.SeverityLevel) {
	t := appinsights.NewTraceTelemetry(name, severity)
	ktc.Track(t)
}
