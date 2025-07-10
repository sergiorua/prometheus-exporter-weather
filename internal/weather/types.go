package weather

import "time"

type WeatherProvider interface {
	GetWeather(city string) (*WeatherData, error)
	GetWeatherByCoords(lat, lon float64) (*WeatherData, error)
}

type WeatherData struct {
	City          string    `json:"city"`
	Country       string    `json:"country"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	Pressure      float64   `json:"pressure"`
	WindSpeed     float64   `json:"wind_speed"`
	WindDirection float64   `json:"wind_direction"`
	Visibility    float64   `json:"visibility"`
	CloudCover    float64   `json:"cloud_cover"`
	Timestamp     time.Time `json:"timestamp"`
	Condition     string    `json:"condition"`
}

type OpenWeatherMapResponse struct {
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
		Pressure float64 `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Visibility int `json:"visibility"`
	Clouds     struct {
		All float64 `json:"all"`
	} `json:"clouds"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}