// Package config implements load, validate, live reload config file
// 

package config

// uses env CONF

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

var (
	ErrConfEnvNotExists error = errors.New("config environment not exists")
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
	Enabled  bool          `yaml:"e_enabled"`  // turns eBPF monitoring on/off
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
	Enabled  bool   `yaml:"p_enabled"`  // turns Prometheus metrics export on/off
	Endpoint string `yaml:"endpoint"` // specifies the metrics HTTP endpoint
	Port     int    `yaml:"port"`     // defines the Prometheus exporter port
}

// Cnf config struct
type Cnf struct{ atomic.Value }
var cnf Cnf


// New return new config instance
//
// Returns:
//   config instance
//
// Example:
//   cnf, err := config.New()
func New() (Cnf, error) {
	p, err := getConfPath()
	if err != nil {
		return cnf, newConfigError("unable to get config file path", err)
	}
	
	f, err := os.Open(p)
    if err != nil {
        return cnf, newConfigError("unable to open config file", err)
    }
    defer f.Close()
	
	decoder := yaml.NewDecoder(f)
	var c Config
    if err := decoder.Decode(&c); err != nil {
        return cnf, newConfigError("unable to decode config", err)
    }
	
	// TODO add validation
	// if err := c.Validate(); err != nil {
    // 	return cnf, newConfigError("invalid config", err)
	// }

	cnf.Store(c)
	return cnf, nil
}


// getConfPath return config path from env CONF, if CONF is not exists return err
//
// Returns:
//   string - config path
//
// Example:
//   p, err := getConfPath()
func getConfPath() (string, error) {
	p, ok := os.LookupEnv(CENV)
	if !ok {
		return p, ErrConfEnvNotExists
	}

	return p, nil
}

// Watch start live reload config file. If context will done - watcher will close
//
// Receiver:
//   Cnf
//
// Parameters:
//   ctx - context
//
// Returns:
//   error - error
//
// Example:
//   ctx := context.Background()
//   defer ctx.Done()
//   err = cnf.Watch(ctx)
func (c *Cnf) Watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return newConfigError("unable to start congif watcher", err)
	}
	
	go func() {
		<-ctx.Done()
		watcher.Close()
	}()
	
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if newCfg, err := New(); err == nil {
					log.Println("⚡️ Конфиг перезагружен:", newCfg)
				} else {
					log.Println("Ошибка при перезагрузке конфига:", err)
				}
			}
			case <-ctx.Done():
				return
			}
		}
	}()
	
	p, err := getConfPath()
	if err != nil {
		watcher.Close()
		return newConfigError("unable to get config file path", err)
	}
	
	if err := watcher.Add(p); err != nil {
		watcher.Close()
		return newConfigError("unable to add config file in config watcher", err)
	}

	return nil
}

// TODO add validation
// func (c *Config) Validate() error {
//     if c.General.HTTPPort == c.General.GRPCPort {
//         return errors.New("HTTP and gRPC ports cannot be the same")
//     }
//     // Add more validation rules
//     return nil
// }
