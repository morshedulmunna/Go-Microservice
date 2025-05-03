package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	App struct {
		Name    string `mapstructure:"name"`
		Env     string `mapstructure:"env"`
		Debug   bool   `mapstructure:"debug"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`

	HTTP struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"http"`

	GRPC struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"grpc"`

	Database struct {
		Driver                 string `mapstructure:"driver"`
		Host                   string `mapstructure:"host"`
		Port                   int    `mapstructure:"port"`
		Name                   string `mapstructure:"name"`
		User                   string `mapstructure:"user"`
		Password               string `mapstructure:"password"`
		SSLMode                string `mapstructure:"sslmode"`
		MaxConnections         int    `mapstructure:"maxConnections"`
		MaxIdleConnections     int    `mapstructure:"maxIdleConnections"`
		MaxLifetimeConnections int    `mapstructure:"maxLifetimeConnections"`
	} `mapstructure:"database"`

	MongoDB struct {
		URI      string `mapstructure:"uri"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongodb"`

	JWT struct {
		SecretKey       string `mapstructure:"secretKey"`
		ExpirationHours int    `mapstructure:"expirationHours"`
	} `mapstructure:"jwt"`

	Consul struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		ServiceName string `mapstructure:"serviceName"`
		ServiceID   string `mapstructure:"serviceID"`
	} `mapstructure:"consul"`

	RabbitMQ struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		VHost    string `mapstructure:"vhost"`
	} `mapstructure:"rabbitmq"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Vault struct {
		Address string `mapstructure:"address"`
		Token   string `mapstructure:"token"`
	} `mapstructure:"vault"`

	Observability struct {
		LogLevel       string `mapstructure:"logLevel"`
		JaegerHost     string `mapstructure:"jaegerHost"`
		JaegerPort     int    `mapstructure:"jaegerPort"`
		PrometheusPort int    `mapstructure:"prometheusPort"`
	} `mapstructure:"observability"`

	RateLimit struct {
		Requests int `mapstructure:"requests"`
		Duration int `mapstructure:"duration"`
	} `mapstructure:"rateLimit"`

	CircuitBreaker struct {
		Timeout     int `mapstructure:"timeout"`
		MaxFailures int `mapstructure:"maxFailures"`
	} `mapstructure:"circuitBreaker"`

	CORS struct {
		AllowedOrigins []string `mapstructure:"allowedOrigins"`
		AllowedMethods []string `mapstructure:"allowedMethods"`
		AllowedHeaders []string `mapstructure:"allowedHeaders"`
		MaxAge         int      `mapstructure:"maxAge"`
	} `mapstructure:"cors"`

	TLS struct {
		CertFile string `mapstructure:"certFile"`
		KeyFile  string `mapstructure:"keyFile"`
	} `mapstructure:"tls"`

	API struct {
		SwaggerHost    string   `mapstructure:"swaggerHost"`
		SwaggerSchemes []string `mapstructure:"swaggerSchemes"`
		SwaggerVersion string   `mapstructure:"swaggerVersion"`
	} `mapstructure:"api"`

	HealthCheck struct {
		Port     int    `mapstructure:"port"`
		Endpoint string `mapstructure:"endpoint"`
	} `mapstructure:"healthCheck"`

	Features struct {
		GRPCEnabled    bool `mapstructure:"grpcEnabled"`
		TracingEnabled bool `mapstructure:"tracingEnabled"`
		MetricsEnabled bool `mapstructure:"metricsEnabled"`
	} `mapstructure:"features"`
}

// Load loads the configuration from config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
