package main

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type serverConfig struct {
	Port        int    `yaml:"listen_port"`
	Url         string `yaml:"public_url"`
	UseAutoCert bool   `yaml:"use_auto_cert"`
	CertPath    string `yaml:"cert_path"`
}

type loggingConfig struct {
	File string `yaml:"file"`
	JSON bool   `yaml:"json"`
	Prod bool   `yaml:"production"`
}

type authConfig struct {
	Password           string `yaml:"password"`
	ConnectionsPerSecond  int    `yaml:"connections_per_second"`
	MaxNamespaceLength int    `yaml:"max_namespace_length"`
}

type config struct {
	Server        serverConfig  `yaml:"server"`
	Logging       loggingConfig `yaml:"logging"`
	Auth          authConfig    `yaml:"auth"`
	logFileHandle io.WriteCloser
}

func getDefaultConfig() *config {
	cfg := &config{}
	cfg.Server.Url = "messaget.example.com"
	cfg.Server.UseAutoCert = false
	cfg.Server.Port = 443
	cfg.Server.CertPath = "/var/www/.cache"
	cfg.Auth.ConnectionsPerSecond = 2
	cfg.Auth.Password = "super-secure-password"
	cfg.Auth.MaxNamespaceLength = 50
	cfg.Logging.Prod = true
	return cfg
}

func (cfg *config) setupLogging() error {
	if cfg.logFileHandle != nil {
		cfg.logFileHandle.Close()
	}
	if cfg.Logging.File != "" {
		logFile, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
		cfg.logFileHandle = logFile
	} else {
		log.SetOutput(os.Stdout)
		cfg.logFileHandle = nil
	}
	if cfg.Logging.JSON {
		log.SetFlags(0)
	} else {
		log.SetFlags(log.LstdFlags)
	}
	return nil
}

func getConfig(configString string, file *string) (*config, error) {
	cfg := getDefaultConfig()

	if err := yaml.UnmarshalStrict([]byte(configString), cfg); err != nil {
		return nil, err
	}

	if err := cfg.setupLogging(); err != nil {
		return nil, err
	}

	// save default if it doesn't exist
	if _, err := os.Stat(*file); os.IsNotExist(err) {
		d, err := yaml.Marshal(&cfg)
		if err != nil {
			errorLogger.Fatalf("Failed to get config: %v", err)
		}
		f, createError := os.Create(*file)
		if createError != nil {
			errorLogger.Fatalf("Failed to create config: %v, %s", err, *file)
		}
		var _, _ = f.WriteString(string(d))
		defer f.Close()
	}

	return cfg, nil
}
