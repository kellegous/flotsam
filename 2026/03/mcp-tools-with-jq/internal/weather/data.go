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
	Time              LocalTime         `json:"time"`
	CurrentConditions CurrentConditions `json:"currently"`
	HourlyForecast    []HourlyForecast  `json:"hourly"`
	DailyForecast     []DailyForecast   `json:"daily"`
}

type CurrentConditions struct {
	Time                 LocalTime `json:"time" jsonschema:"The time of the current conditions"`
	Summary              string    `json:"summary" jsonschema:"A brief description of the current weather conditions"`
	Icon                 string    `json:"icon" jsonschema:"The icon representing the summary"`
	Temperature          float64   `json:"temperature" jsonschema:"The current temperature in degrees Fahrenheit"`
	ApparentTemperature  float64   `json:"apparent_temperature" jsonschema:"The apparent temperature, taking into account wind and humdity"`
	DewPoint             float64   `json:"dew_point" jsonschema:"The dew point in degrees Fahrenheit"`
	Humidity             float64   `json:"humidity" jsonschema:"The current humidity as a percentage"`
	WindSpeed            float64   `json:"wind_speed" jsonschema:"The wind speed in miles per hour"`
	WindBearing          float64   `json:"wind_bearing" jsonschema:"The wind bearing in degrees"`
	Visibility           float64   `json:"visibility" jsonschema:"The current visibility in miles"`
	CloudCover           float64   `json:"cloud_cover" jsonschema:"The current cloud cover as a percentage"`
	Pressure             float64   `json:"pressure" jsonschema:"The current pressure in millibars"`
	NearestStormBearing  float64   `json:"nearest_storm_bearing" jsonschema:"The bearing in degrees to the nearest storm from the current location"`
	NearestStormDistance float64   `json:"nearest_storm_distance" jsonschema:"The distance in miles to the nearest storm from the current location"`
}

type HourlyForecast struct {
	Time                LocalTime `json:"time" jsonschema:"The time of the top of the hour for this forecast"`
	Summary             string    `json:"summary" jsonschema:"A brief description of the weather conditions"`
	Icon                string    `json:"icon" jsonschema:"The icon representing the summary"`
	Temperature         float64   `json:"temperature" jsonschema:"The temperature in degrees Fahrenheit"`
	ApparentTemperature float64   `json:"apparent_temperature" jsonschema:"The apparent temperature, taking into account wind and humdity"`
	DewPoint            float64   `json:"dew_point" jsonschema:"The dew point in degrees Fahrenheit"`
	Humidity            float64   `json:"humidity" jsonschema:"The current humidity as a percentage"`
	WindSpeed           float64   `json:"wind_speed" jsonschema:"The wind speed in miles per hour"`
	WindBearing         float64   `json:"wind_bearing" jsonschema:"The wind bearing in degrees"`
	Visibility          float64   `json:"visibility" jsonschema:"The visibility in miles"`
	CloudCover          float64   `json:"cloud_cover" jsonschema:"The cloud cover as a percentage"`
	Pressure            float64   `json:"pressure" jsonschema:"The pressure in millibars"`
}

type DailyForecast struct {
	Time                       LocalTime `json:"time" jsonschema:"The time of the beginning of the day for this forecast"`
	Summary                    string    `json:"summary" jsonschema:"A brief description of the weather conditions"`
	Icon                       string    `json:"icon" jsonschema:"The icon representing the summary"`
	SunriseTime                LocalTime `json:"sunrise_time" jsonschema:"The time of the sunrise"`
	SunsetTime                 LocalTime `json:"sunset_time" jsonschema:"The time of the sunset"`
	MoonPhase                  float64   `json:"moon_phase" jsonschema:"The phase of the moon"`
	TemperatureMin             float64   `json:"temperature_min" jsonschema:"The minimum temperature in degrees Fahrenheit"`
	TemperatureMinTime         LocalTime `json:"temperature_min_time" jsonschema:"The time of the minimum temperature"`
	TemperatureMax             float64   `json:"temperature_max" jsonschema:"The maximum temperature in degrees Fahrenheit"`
	TemperatureMaxTime         time.Time `json:"temperature_max_time" jsonschema:"The time of the maximum temperature"`
	ApparentTemperatureMax     float64   `json:"apparent_temperature_max" jsonschema:"The apparent temperature, taking into account wind and humdity"`
	ApparentTemperatureMinTime LocalTime `json:"apparent_temperature_min_time" jsonschema:"The time of the minimum apparent temperature"`
	DewPoint                   float64   `json:"dew_point" jsonschema:"The dew point in degrees Fahrenheit"`
	Humidity                   float64   `json:"humidity" jsonschema:"The current humidity as a percentage"`
	WindSpeed                  float64   `json:"wind_speed" jsonschema:"The wind speed in miles per hour"`
	WindBearing                float64   `json:"wind_bearing" jsonschema:"The wind bearing in degrees"`
	Visibility                 float64   `json:"visibility" jsonschema:"The visibility in miles"`
	CloudCover                 float64   `json:"cloud_cover" jsonschema:"The cloud cover as a percentage"`
	Pressure                   float64   `json:"pressure" jsonschema:"The pressure in millibars"`
	PrecipitationProbability   float64   `json:"precipitation_probability,omitempty" jsonschema:"The probability of precipitation as a percentage"`
	PrecipitationIntensity     float64   `json:"precipitation_intensity,omitempty" jsonschema:"The precipitation intensity in inches per hour"`
	PrecipitationIntensityMax  float64   `json:"precipitation_intensity_max,omitempty" jsonschema:"The maximum precipitation intensity in inches per hour"`
}

func Get() (*Data, error) {
	var d Data
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, poop.Chain(err)
	}
	return &d, nil
}
