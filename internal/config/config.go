package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Weather  WeatherConfig  `mapstructure:"weather"`
	Cities   []CityConfig   `mapstructure:"cities"`
	Scraping ScrapingConfig `mapstructure:"scraping"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port" envconfig:"SERVER_PORT" default:"8080"`
	MetricsPath  string        `mapstructure:"metrics_path" envconfig:"METRICS_PATH" default:"/metrics"`
	HealthPath   string        `mapstructure:"health_path" envconfig:"HEALTH_PATH" default:"/health"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" envconfig:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" envconfig:"WRITE_TIMEOUT" default:"30s"`
}

type WeatherConfig struct {
	APIKey  string `mapstructure:"api_key" envconfig:"OPENWEATHER_API_KEY" required:"true"`
	BaseURL string `mapstructure:"base_url" envconfig:"WEATHER_BASE_URL" default:"https://api.openweathermap.org/data/2.5"`
	Timeout time.Duration `mapstructure:"timeout" envconfig:"WEATHER_TIMEOUT" default:"10s"`
}

type CityConfig struct {
	Name    string      `mapstructure:"name"`
	Country string      `mapstructure:"country"`
	Coords  Coordinates `mapstructure:"coordinates"`
}

type Coordinates struct {
	Lat float64 `mapstructure:"lat"`
	Lon float64 `mapstructure:"lon"`
}

type ScrapingConfig struct {
	Interval      time.Duration `mapstructure:"interval" envconfig:"SCRAPING_INTERVAL" default:"300s"`
	Timeout       time.Duration `mapstructure:"timeout" envconfig:"SCRAPING_TIMEOUT" default:"30s"`
	RetryAttempts int           `mapstructure:"retry_attempts" envconfig:"RETRY_ATTEMPTS" default:"3"`
	RetryDelay    time.Duration `mapstructure:"retry_delay" envconfig:"RETRY_DELAY" default:"10s"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level" envconfig:"LOG_LEVEL" default:"info"`
	Format string `mapstructure:"format" envconfig:"LOG_FORMAT" default:"json"`
	Output string `mapstructure:"output" envconfig:"LOG_OUTPUT" default:"stdout"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config
	
	// Set default cities if none specified
	cfg.Cities = []CityConfig{
		{Name: "London", Country: "UK", Coords: Coordinates{Lat: 51.5074, Lon: -0.1278}},
		{Name: "New York", Country: "US", Coords: Coordinates{Lat: 40.7128, Lon: -74.0060}},
		{Name: "Tokyo", Country: "JP", Coords: Coordinates{Lat: 35.6762, Lon: 139.6503}},
	}
	
	// Load from file if exists
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		} else {
			if err := viper.Unmarshal(&cfg); err != nil {
				return nil, fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}
	}
	
	// Override with environment variables
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}
	
	// Validate required fields
	if cfg.Weather.APIKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY is required")
	}
	
	return &cfg, nil
}