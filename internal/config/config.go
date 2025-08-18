package config

// uses env CONF

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfEnvNotExists error = errors.New("config envinronment not exists")
)

const (
	CENV = "CONF_PATH"
)

// ConfigError struct fot custom errors
type ConfigError struct {
	msg string
	err error
}

// Error satisfy the error interface condition
func (c *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s: %v", c.msg, c.err)
}

// newConfigError attach err to msg and return new error
func newConfigError(msg string, err error) error {
	return &ConfigError{msg: msg, err: err}
}

// Config represents the root configuration structure.
// For development, configs/config.yaml is used.
// For production in Docker, /etc/diskmon/config.yaml is mounted.
type Config struct {
	General    General    `yaml:"general"`    // common application settings
	Monitoring Monitoring `yaml:"monitoring"` // system monitoring configuration
	EBPF       EBPF       `yaml:"ebpf"`       // eBPF monitoring configuration
	Cgroups    Cgroups    `yaml:"cgroups"`    // control groups configuration
	Alerting   Alerting   `yaml:"alerting"`   // alert notification settings
	Prometheus Prometheus `yaml:"prometheus"` // metrics export settings
}

// General contains common application settings
type General struct {
	LogLevel string `yaml:"log_level"` // defines the logging verbosity (debug/info/warn/error)
	HTTPPort int    `yaml:"http_port"` // specifies the port for REST API server
	GRPCPort int    `yaml:"grpc_port"` // specifies the port for gRPC server
}

// Monitoring contains system monitoring configuration
type Monitoring struct {
	Interval string `yaml:"interval"`  // defines metrics collection frequency (e.g., "5s")
	ProcPath string `yaml:"proc_path"` //specifies the path to proc filesystem (for testing)
}

// EBPF contains eBPF monitoring configuration
type EBPF struct {
	Enabled  bool          `yaml:"enabled"`  // turns eBPF monitoring on/off
	Programs []EBPFProgram `yaml:"programs"` //  eBPF programs to load
}

// EBPFProgram defines a single eBPF program configuration
type EBPFProgram struct {
	Name      string `yaml:"name"`       // identifies the eBPF program
	ProbeType string `yaml:"probe_type"` // specifies the attachment type (kprobe/tracepoint)
	Target    string `yaml:"target"`     // specifies the system call to monitor
}

// Cgroups contains control groups configuration
type Cgroups struct {
	BasePath      string `yaml:"base_path"`      // specifies the path to cgroups v2 filesystem
	DefaultWeight int    `yaml:"default_weight"` // defines the default I/O weight (1-100)
}

// Alerting contains alert notification settings
type Alerting struct {
	DiskLatencyThresholdMS int         `yaml:"disk_latency_threshold_ms"` // defines alert threshold in milliseconds
	IOPSThreshold          int         `yaml:"iops_threshold"`            // defines maximum IO operations per second
	NotifyMethods          []string    `yaml:"notify_methods"`            // lists enabled notification methods
	Slack                  SlackConfig `yaml:"slack"`                     // contains Slack integration settings
	Email                  EmailConfig `yaml:"email"`                     // contains email notification settings
}

// SlackConfig contains Slack webhook settings
type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"` // specifies the Slack incoming webhook URL
	Channel    string `yaml:"channel"`     // defines the target Slack channel for alerts
}

// EmailConfig contains SMTP notification settings
type EmailConfig struct {
	SMTPServer string `yaml:"smtp_server"` // specifies the SMTP server address
	From       string `yaml:"from"`        // defines the sender email address
	To         string `yaml:"to"`          // specifies the recipient email address
}

// Prometheus contains metrics export settings
type Prometheus struct {
	Enabled  bool   `yaml:"enabled"`  // turns Prometheus metrics export on/off
	Endpoint string `yaml:"endpoint"` // specifies the metrics HTTP endpoint
	Port     int    `yaml:"port"`     // defines the Prometheus exporter port
}

// New return new config instance
func New() (*Config, error) {
	var c Config

	p, err := getConfPath()
	if err != nil {
		return nil, newConfigError("unable to get config file path", err)
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return nil, newConfigError("unable to read config file", err)
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, newConfigError("unable to unmarshall config file", err)

	}

	return &c, nil
}

// getConfPath return config path from env CONF, if CONF is not exists return err
func getConfPath() (string, error) {
	p, ok := os.LookupEnv(CENV)
	if !ok {
		return p, ErrConfEnvNotExists
	}

	return p, nil
}
