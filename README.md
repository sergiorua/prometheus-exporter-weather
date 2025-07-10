# Weather Exporter

A Prometheus exporter that fetches weather data from OpenWeatherMap API and exposes it as metrics.

## Features

- Fetches weather data for multiple cities
- Exposes metrics in Prometheus format
- Configurable through YAML files and environment variables
- Docker support
- Health and readiness endpoints

## Quick Start

### Prerequisites

- Go 1.21+
- OpenWeatherMap API key (get one at https://openweathermap.org/api)

### Running Locally

1. Set your API key:
   ```bash
   export OPENWEATHER_API_KEY="your-api-key-here"
   ```

2. Build and run:
   ```bash
   make run
   ```

3. Access metrics at http://localhost:8080/metrics

### Running with Docker

```bash
make docker-run
```

## Configuration

The application can be configured via:
- Configuration file (`configs/config.yaml`)
- Environment variables

### Environment Variables

- `OPENWEATHER_API_KEY` (required): Your OpenWeatherMap API key
- `SERVER_PORT`: HTTP server port (default: 8080)
- `SCRAPING_INTERVAL`: Weather data collection interval (default: 300s)
- `LOG_LEVEL`: Logging level (default: info)

## Metrics

The exporter provides the following metrics:

- `weather_temperature_celsius`: Current temperature in Celsius
- `weather_humidity_percent`: Current humidity percentage
- `weather_pressure_hpa`: Atmospheric pressure in hPa
- `weather_wind_speed_mps`: Wind speed in meters per second
- `weather_wind_direction_degrees`: Wind direction in degrees
- `weather_visibility_km`: Visibility in kilometers
- `weather_cloud_cover_percent`: Cloud cover percentage
- `weather_api_requests_total`: Total API requests counter
- `weather_api_request_duration_seconds`: API request duration histogram
- `weather_last_update_timestamp`: Last successful update timestamp

## Endpoints

- `/metrics` - Prometheus metrics
- `/health` - Health check
- `/ready` - Readiness check
- `/weather/{city}` - Current weather for a specific city (JSON)

## Building

```bash
# Build binary
make build

# Run tests
make test

# Build Docker image
make docker-build
```