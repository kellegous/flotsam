package weather

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/kellegous/poop"
)

//go:embed weather.json
var data []byte

type Data struct {
	Time              time.Time         `json:"time"`
	CurrentConditions CurrentConditions `json:"currently"`
	HourlyForecast    []HourlyForecast  `json:"hourly"`
	DailyForecast     []DailyForecast   `json:"daily"`
}

type CurrentConditions struct {
	Time                 time.Time `json:"time"`
	Summary              string    `json:"summary"`
	Icon                 string    `json:"icon"`
	Temperature          float64   `json:"temperature"`
	ApparentTemperature  float64   `json:"apparent_temperature"`
	DewPoint             float64   `json:"dew_point"`
	Humidity             float64   `json:"humidity"`
	WindSpeed            float64   `json:"wind_speed"`
	WindBearing          float64   `json:"wind_bearing"`
	Visibility           float64   `json:"visibility"`
	CloudCover           float64   `json:"cloud_cover"`
	Pressure             float64   `json:"pressure"`
	NearestStormBearing  float64   `json:"nearest_storm_bearing"`
	NearestStormDistance float64   `json:"nearest_storm_distance"`
}

type HourlyForecast struct {
	Time                time.Time `json:"time"`
	Summary             string    `json:"summary"`
	Icon                string    `json:"icon"`
	Temperature         float64   `json:"temperature"`
	ApparentTemperature float64   `json:"apparent_temperature"`
	DewPoint            float64   `json:"dew_point"`
	Humidity            float64   `json:"humidity"`
	WindSpeed           float64   `json:"wind_speed"`
	WindBearing         float64   `json:"wind_bearing"`
	Visibility          float64   `json:"visibility"`
	CloudCover          float64   `json:"cloud_cover"`
	Pressure            float64   `json:"pressure"`
}

type DailyForecast struct {
	Time                       time.Time `json:"time"`
	Summary                    string    `json:"summary"`
	Icon                       string    `json:"icon"`
	SunriseTime                time.Time `json:"sunrise_time"`
	SunsetTime                 time.Time `json:"sunset_time"`
	MoonPhase                  float64   `json:"moon_phase"`
	TemperatureMin             float64   `json:"temperature_min"`
	TemperatureMinTime         time.Time `json:"temperature_min_time"`
	TemperatureMax             float64   `json:"temperature_max"`
	TemperatureMaxTime         time.Time `json:"temperature_max_time"`
	ApparentTemperatureMax     float64   `json:"apparent_temperature_max"`
	ApparentTemperatureMinTime time.Time `json:"apparent_temperature_min_time"`
	DewPoint                   float64   `json:"dew_point"`
	Humidity                   float64   `json:"humidity"`
	WindSpeed                  float64   `json:"wind_speed"`
	WindBearing                float64   `json:"wind_bearing"`
	Visibility                 float64   `json:"visibility"`
	CloudCover                 float64   `json:"cloud_cover"`
	Pressure                   float64   `json:"pressure"`
	PrecipitationProbability   float64   `json:"precipitation_probability,omitempty"`
	PrecipitationIntensity     float64   `json:"precipitation_intensity,omitempty"`
	PrecipitationIntensityMax  float64   `json:"precipitation_intensity_max,omitempty"`
}

func Get() (*Data, error) {
	var d Data
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, poop.Chain(err)
	}
	return &d, nil
}
