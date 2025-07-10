package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type OpenWeatherMapClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewOpenWeatherMapClient(apiKey string) *OpenWeatherMapClient {
	return &OpenWeatherMapClient{
		apiKey:  apiKey,
		baseURL: "https://api.openweathermap.org/data/2.5",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *OpenWeatherMapClient) GetWeather(city string) (*WeatherData, error) {
	encodedCity := url.QueryEscape(city)
	apiURL := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=metric", c.baseURL, encodedCity, c.apiKey)
	
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	var owmResp OpenWeatherMapResponse
	if err := json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return c.convertToWeatherData(&owmResp), nil
}

func (c *OpenWeatherMapClient) GetWeatherByCoords(lat, lon float64) (*WeatherData, error) {
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=metric", c.baseURL, lat, lon, c.apiKey)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}
	
	var owmResp OpenWeatherMapResponse
	if err := json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return c.convertToWeatherData(&owmResp), nil
}

func (c *OpenWeatherMapClient) convertToWeatherData(resp *OpenWeatherMapResponse) *WeatherData {
	condition := "Unknown"
	if len(resp.Weather) > 0 {
		condition = resp.Weather[0].Main
	}
	
	return &WeatherData{
		City:          resp.Name,
		Country:       resp.Sys.Country,
		Temperature:   resp.Main.Temp,
		Humidity:      resp.Main.Humidity,
		Pressure:      resp.Main.Pressure,
		WindSpeed:     resp.Wind.Speed,
		WindDirection: resp.Wind.Deg,
		Visibility:    float64(resp.Visibility) / 1000.0,
		CloudCover:    resp.Clouds.All,
		Timestamp:     time.Now(),
		Condition:     condition,
	}
}