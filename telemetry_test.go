package appink8s

import (
	"errors"
	"testing"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"

	"github.com/stretchr/testify/assert"
)

type telemetry_mockTelemetryClient struct {
	ctx     *appinsights.TelemetryContext
	tracked appinsights.Telemetry
}

func (m *telemetry_mockTelemetryClient) Channel() appinsights.TelemetryChannel {
	return nil
}
func (m *telemetry_mockTelemetryClient) Context() *appinsights.TelemetryContext {
	return m.ctx
}
func (m *telemetry_mockTelemetryClient) InstrumentationKey() string {
	return ""
}
func (m *telemetry_mockTelemetryClient) IsEnabled() bool {
	return true
}
func (_ *telemetry_mockTelemetryClient) SetIsEnabled(_ bool) {}
func (m *telemetry_mockTelemetryClient) Track(t appinsights.Telemetry) {
	m.tracked = t
}
func (_ *telemetry_mockTelemetryClient) TrackAvailability(name string, duration time.Duration, success bool) {
}
func (_ *telemetry_mockTelemetryClient) TrackEvent(name string)                 {}
func (_ *telemetry_mockTelemetryClient) TrackException(err interface{})         {}
func (_ *telemetry_mockTelemetryClient) TrackMetric(name string, value float64) {}
func (_ *telemetry_mockTelemetryClient) TrackRemoteDependency(name, dependencyType, target string, success bool) {
}
func (_ *telemetry_mockTelemetryClient) TrackRequest(method, uri string, duration time.Duration, responseCode string) {
}
func (_ *telemetry_mockTelemetryClient) TrackTrace(name string, severity contracts.SeverityLevel) {}

type telemetry_mockInitializer struct {
	called int
	spec   *runtimeSpec
	err    error
}

func (m *telemetry_mockInitializer) ReadPropertySpec() (*runtimeSpec, error) {
	m.called = m.called + 1

	if m.spec != nil {
		return m.spec, m.err
	}

	return &runtimeSpec{}, m.err
}

func Test_That_Apply_Initializes_Property_Handling_When_Uninitialized(t *testing.T) {
	c := &kubernetesTelemetryClient{
		active:      true,
		initialized: false,
		initializer: &telemetry_mockInitializer{},
	}

	m := make(map[string]string)
	c.apply(m)

	assert.True(t, c.initialized)
}

func Test_That_Apply_Initializes_Property_Handling_Only_Once(t *testing.T) {
	i := &telemetry_mockInitializer{}
	c := &kubernetesTelemetryClient{
		active:      true,
		initialized: false,
		initializer: i,
	}

	m := make(map[string]string)
	c.apply(m)
	c.apply(m)
	c.apply(m)

	assert.Equal(t, 1, i.called)
}

func Test_That_Apply_Deactivates_Telemetry_Enhancements_On_Initialization_Error(t *testing.T) {
	c := &kubernetesTelemetryClient{
		initialized: false,
		initializer: &telemetry_mockInitializer{err: errors.New("mock")},
	}

	m := make(map[string]string)
	c.apply(m)

	assert.False(t, c.active)
}

func Test_That_Apply_Adds_Telemetry_Enhancing_Properties_To_Property_Map(t *testing.T) {
	p := newSpec().ToPropertyMap()

	c := &kubernetesTelemetryClient{
		active:      true,
		initialized: true,
		properties:  p,
	}

	m := make(map[string]string)
	c.apply(m)

	assert.Equal(t, p, m)
}

func Test_That_Apply_Skips_Telemetry_Enhancing_Properties_When_Deactivated(t *testing.T) {
	s := newSpec()
	p := s.ToPropertyMap()

	c := &kubernetesTelemetryClient{
		active:      false,
		initialized: true,
		initializer: &telemetry_mockInitializer{spec: s},
		properties:  p,
	}

	m := make(map[string]string)
	c.apply(m)

	assert.NotEqual(t, p, m)
}

func Test_That_Initialize_Assigns_Telemetry_Context_Role_To_DeploymentName(t *testing.T) {
	s := newSpec()

	c := &kubernetesTelemetryClient{
		TelemetryClient: &telemetry_mockTelemetryClient{
			ctx: appinsights.NewTelemetryContext(""),
		},
		active:      true,
		initialized: false,
		initializer: &telemetry_mockInitializer{spec: s},
		properties:  s.ToPropertyMap(),
	}

	c.initialize()

	role := c.TelemetryClient.Context().Tags.Cloud().GetRole()
	assert.Equal(t, s.DeploymentName, role)
}

func Test_That_Initialize_Assigns_Telemetry_Context_RoleInstance_To_ReplicaSetName(t *testing.T) {
	s := newSpec()

	c := &kubernetesTelemetryClient{
		TelemetryClient: &telemetry_mockTelemetryClient{
			ctx: appinsights.NewTelemetryContext(""),
		},
		active:      true,
		initialized: false,
		initializer: &telemetry_mockInitializer{spec: s},
		properties:  s.ToPropertyMap(),
	}

	c.initialize()

	instance := c.TelemetryClient.Context().Tags.Cloud().GetRoleInstance()
	assert.Equal(t, s.ReplicaSetName, instance)
}

func Test_That_Track_Adds_Kubernetes_Properties_To_Telemetry(t *testing.T) {
	s := newSpec()
	p := s.ToPropertyMap()

	c := &kubernetesTelemetryClient{
		TelemetryClient: &telemetry_mockTelemetryClient{
			ctx: appinsights.NewTelemetryContext(""),
		},
		active:      true,
		initialized: false,
		initializer: &telemetry_mockInitializer{spec: s},
		properties:  p,
	}

	m := appinsights.NewEventTelemetry("test")
	c.Track(m)

	assert.Equal(t, p, m.GetProperties())
}

func newSpec() *runtimeSpec {
	return &runtimeSpec{
		ContainerID:    "container-id",
		ContainerName:  "container-name",
		DeploymentName: "deployment-name",
		NodeID:         "node-id",
		NodeLabels:     "node-labels",
		NodeName:       "node-name",
		PodID:          "pod-id",
		PodLabels:      "pod-labels",
		PodName:        "pod-name",
		ReplicaSetName: "replicaset-name",
	}
}
