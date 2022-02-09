package config

// Constants
const prometheusTag string = "prom/prometheus:v2.31.1"

// Defaults
const defaultPrometheusPort uint16 = 9091
const defaultPrometheusOpenPort bool = false

// Configuration for Prometheus
type PrometheusConfig struct {
	// The master configuration this belongs to
	MasterConfig *MasterConfig

	// The port to serve metrics on
	Port *Parameter

	// Toggle for forwarding the API port outside of Docker
	OpenPort *Parameter

	// The Docker Hub tag for Prometheus
	ContainerTag *Parameter
}

// Generates a new Prometheus config
func NewPrometheusConfig(config *MasterConfig) *PrometheusConfig {
	return &PrometheusConfig{
		MasterConfig: config,

		Port: &Parameter{
			ID:                   "port",
			Name:                 "API Port",
			Description:          "The port Prometheus should make its statistics available on.",
			Type:                 ParameterType_Uint16,
			Default:              map[Network]interface{}{Network_All: defaultPrometheusPort},
			AffectsContainers:    []ContainerID{ContainerID_Prometheus},
			EnvironmentVariables: []string{"PROMETHEUS_PORT"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},

		OpenPort: &Parameter{
			ID:                   "openPort",
			Name:                 "Open Port",
			Description:          "Enable this to open Prometheus's port to your local network, so other machines can access it too.",
			Type:                 ParameterType_Bool,
			Default:              map[Network]interface{}{Network_All: defaultPrometheusOpenPort},
			AffectsContainers:    []ContainerID{ContainerID_Prometheus},
			EnvironmentVariables: []string{"PROMETHEUS_OPEN_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		ContainerTag: &Parameter{
			ID:                   "containerTag",
			Name:                 "Container Tag",
			Description:          "The tag name of the Prometheus container you want to use on Docker Hub.",
			Type:                 ParameterType_String,
			Default:              map[Network]interface{}{Network_All: prometheusTag},
			AffectsContainers:    []ContainerID{ContainerID_Prometheus},
			EnvironmentVariables: []string{"PROMETHEUS_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},
	}
}
