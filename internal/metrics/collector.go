package metrics

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/weather-exporter/internal/weather"
)

var (
	temperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_temperature_celsius",
			Help: "Current temperature in Celsius",
		},
		[]string{"city", "country", "provider"},
	)

	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_humidity_percent",
			Help: "Current humidity percentage",
		},
		[]string{"city", "country", "provider"},
	)

	pressureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_pressure_hpa",
			Help: "Current atmospheric pressure in hPa",
		},
		[]string{"city", "country", "provider"},
	)

	windSpeedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_wind_speed_mps",
			Help: "Current wind speed in meters per second",
		},
		[]string{"city", "country", "provider"},
	)

	windDirectionGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_wind_direction_degrees",
			Help: "Current wind direction in degrees",
		},
		[]string{"city", "country", "provider"},
	)

	visibilityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_visibility_km",
			Help: "Current visibility in kilometers",
		},
		[]string{"city", "country", "provider"},
	)

	cloudCoverGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_cloud_cover_percent",
			Help: "Current cloud cover percentage",
		},
		[]string{"city", "country", "provider"},
	)

	apiRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_api_requests_total",
			Help: "Total number of API requests to weather providers",
		},
		[]string{"provider", "status"},
	)

	apiRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "weather_api_request_duration_seconds",
			Help:    "Duration of weather API requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"provider"},
	)

	lastUpdateTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "weather_last_update_timestamp",
			Help: "Timestamp of last successful weather data update",
		},
		[]string{"city", "country", "provider"},
	)
)

func init() {
	prometheus.MustRegister(
		temperatureGauge,
		humidityGauge,
		pressureGauge,
		windSpeedGauge,
		windDirectionGauge,
		visibilityGauge,
		cloudCoverGauge,
		apiRequestsTotal,
		apiRequestDuration,
		lastUpdateTimestamp,
	)
}

type Collector struct {
	provider weather.WeatherProvider
	cities   []string
	interval time.Duration
	mutex    sync.RWMutex
}

func NewCollector(provider weather.WeatherProvider, cities []string, interval time.Duration) *Collector {
	return &Collector{
		provider: provider,
		cities:   cities,
		interval: interval,
	}
}

func (c *Collector) Start(ctx context.Context) {
	c.collectMetrics()
	
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.collectMetrics()
		case <-ctx.Done():
			return
		}
	}
}

func (c *Collector) collectMetrics() {
	var wg sync.WaitGroup
	
	for _, city := range c.cities {
		wg.Add(1)
		go func(cityName string) {
			defer wg.Done()
			c.collectCityWeather(cityName)
		}(city)
	}
	
	wg.Wait()
}

func (c *Collector) collectCityWeather(city string) {
	start := time.Now()
	providerName := "openweathermap"
	
	weatherData, err := c.provider.GetWeather(city)
	duration := time.Since(start).Seconds()
	
	apiRequestDuration.WithLabelValues(providerName).Observe(duration)
	
	if err != nil {
		log.Printf("Failed to fetch weather for %s: %v", city, err)
		apiRequestsTotal.WithLabelValues(providerName, "error").Inc()
		return
	}
	
	apiRequestsTotal.WithLabelValues(providerName, "success").Inc()
	
	labels := prometheus.Labels{
		"city":     weatherData.City,
		"country":  weatherData.Country,
		"provider": providerName,
	}
	
	temperatureGauge.With(labels).Set(weatherData.Temperature)
	humidityGauge.With(labels).Set(weatherData.Humidity)
	pressureGauge.With(labels).Set(weatherData.Pressure)
	windSpeedGauge.With(labels).Set(weatherData.WindSpeed)
	windDirectionGauge.With(labels).Set(weatherData.WindDirection)
	visibilityGauge.With(labels).Set(weatherData.Visibility)
	cloudCoverGauge.With(labels).Set(weatherData.CloudCover)
	lastUpdateTimestamp.With(labels).Set(float64(weatherData.Timestamp.Unix()))
}