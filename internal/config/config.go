package config

// uses env CONF

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfEnvNotExists error = errors.New("config envinronment not exists")
)

const (
	CENV = "CONF_PATH"
)

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
	// LogLevel defines the logging verbosity (debug/info/warn/error)
	LogLevel string `yaml:"log_level"`

	// HTTPPort specifies the port for REST API server
	HTTPPort int `yaml:"http_port"`

	// GRPCPort specifies the port for gRPC server
	GRPCPort int `yaml:"grpc_port"`
}

// Monitoring contains system monitoring configuration
type Monitoring struct {
	// Interval defines metrics collection frequency (e.g., "5s")
	Interval string `yaml:"interval"`

	// ProcPath specifies the path to proc filesystem (for testing)
	ProcPath string `yaml:"proc_path"`
}

// EBPF contains eBPF monitoring configuration
type EBPF struct {
	// Enabled turns eBPF monitoring on/off
	Enabled bool `yaml:"enabled"`

	// Programs lists eBPF programs to load
	Programs []EBPFProgram `yaml:"programs"`
}

// EBPFProgram defines a single eBPF program configuration
type EBPFProgram struct {
	// Name identifies the eBPF program
	Name string `yaml:"name"`

	// ProbeType specifies the attachment type (kprobe/tracepoint)
	ProbeType string `yaml:"probe_type"`

	// Target specifies the system call to monitor
	Target string `yaml:"target"`
}

// Cgroups contains control groups configuration
type Cgroups struct {
	// BasePath specifies the path to cgroups v2 filesystem
	BasePath string `yaml:"base_path"`

	// DefaultWeight defines the default I/O weight (1-100)
	DefaultWeight int `yaml:"default_weight"`
}

// Alerting contains alert notification settings
type Alerting struct {
	// DiskLatencyThresholdMS defines alert threshold in milliseconds
	DiskLatencyThresholdMS int `yaml:"disk_latency_threshold_ms"`

	// IOPSThreshold defines maximum IO operations per second
	IOPSThreshold int `yaml:"iops_threshold"`

	// NotifyMethods lists enabled notification methods
	NotifyMethods []string `yaml:"notify_methods"`

	// Slack contains Slack integration settings
	Slack SlackConfig `yaml:"slack"`

	// Email contains email notification settings
	Email EmailConfig `yaml:"email"`
}

// SlackConfig contains Slack webhook settings
type SlackConfig struct {
	// WebhookURL specifies the Slack incoming webhook URL
	WebhookURL string `yaml:"webhook_url"`

	// Channel defines the target Slack channel for alerts
	Channel string `yaml:"channel"`
}

// EmailConfig contains SMTP notification settings
type EmailConfig struct {
	// SMTPServer specifies the SMTP server address
	SMTPServer string `yaml:"smtp_server"`

	// From defines the sender email address
	From string `yaml:"from"`

	// To specifies the recipient email address
	To string `yaml:"to"`
}

// Prometheus contains metrics export settings
type Prometheus struct {
	// Enabled turns Prometheus metrics export on/off
	Enabled bool `yaml:"enabled"`

	// Endpoint specifies the metrics HTTP endpoint
	Endpoint string `yaml:"endpoint"`

	// Port defines the Prometheus exporter port
	Port int `yaml:"port"`
}

// New return new config instance
func New() (*Config, error) {
	var c Config
	
	p, err := getConfPath()
	if err != nil {
		return nil, fmt.Errorf("unable to get config file path: %w", err)
	}
	
	f, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}
	
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall config file: %w", err)
	}
	
	return &c, nil
}

// getConfPath return config path from env CONF, if CONF is not exists return err
func getConfPath() (string, error) {
	p, ok := os.LookupEnv(CENV)
	if ok {
		return p, nil
	} else {
		return p, ErrConfEnvNotExists
	}
}
