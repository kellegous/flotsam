package plain

import (
	"context"
	"kellegous/jqmcp/internal/weather"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New(ctx context.Context) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "Local Weather",
		Version: "1.0.0",
	}, nil)

	// s.AddReceivingMiddleware(withLogging(logging.L(ctx)))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "current-time",
		Description: "Get the current time",
	}, currentTime)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "current-conditions",
		Description: "Get the current weather conditions",
	}, getCurrentConditions)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "hourly-forecast",
		Description: "Get the hourly weather forecast for the local area for the next 48 hours",
	}, hourlyForecast)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "daily-forecast",
		Description: "Get the daily weather forecast for the local area for the next 7 days",
	}, dailyForecast)

	return s
}

type currentTimeReq struct{}

type currentTimeRes struct {
	Time string `json:"time" jsonschema:"Current time in local time zone in ISO 8601 format"`
}

func currentTime(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input currentTimeReq,
) (*mcp.CallToolResult, currentTimeRes, error) {
	return nil, currentTimeRes{
		Time: time.Now().Format(time.RFC3339),
	}, nil
}

type getCurrentConditionsReq struct{}

type getCurrentConditionsRes struct {
	Conditions *weather.CurrentConditions `json:"conditions" jsonschema:"Current weather conditions in the local area"`
}

func getCurrentConditions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args getCurrentConditionsReq,
) (*mcp.CallToolResult, getCurrentConditionsRes, error) {
	data, err := weather.Get()
	if err != nil {
		return nil, getCurrentConditionsRes{}, err
	}
	return nil, getCurrentConditionsRes{
		Conditions: &data.CurrentConditions,
	}, nil
}

type hourlyForecastReq struct{}

type hourlyForecastRes struct {
	Forecast []weather.HourlyForecast `json:"forecast" jsonschema:"Hourly weather forecast for the local area for the next 48 hours"`
}

func hourlyForecast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args hourlyForecastReq,
) (*mcp.CallToolResult, hourlyForecastRes, error) {
	data, err := weather.Get()
	if err != nil {
		return nil, hourlyForecastRes{}, err
	}
	return nil, hourlyForecastRes{
		Forecast: data.HourlyForecast,
	}, nil
}

type dailyForecastReq struct{}

type dailyForecastRes struct {
	Forecast []weather.DailyForecast `json:"forecast" jsonschema:"Daily weather forecast for the local area for the next 7 days"`
}

func dailyForecast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args dailyForecastReq,
) (*mcp.CallToolResult, dailyForecastRes, error) {
	data, err := weather.Get()
	if err != nil {
		return nil, dailyForecastRes{}, err
	}
	return nil, dailyForecastRes{
		Forecast: data.DailyForecast,
	}, nil
}
