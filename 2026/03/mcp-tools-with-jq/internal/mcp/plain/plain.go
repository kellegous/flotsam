package plain

import (
	"context"
	"kellegous/jqmcp/internal/mcp/util"
	"kellegous/jqmcp/internal/weather"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New(ctx context.Context) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "Local Weather",
		Version: "1.0.0",
	}, nil)

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
	CurrentTime *util.Time `json:"current_time,omitempty" jsonschema:"Data about the current time in the local time zone"`
}

func currentTime(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input currentTimeReq,
) (*mcp.CallToolResult, currentTimeRes, error) {
	data, err := weather.Get()
	if err != nil {
		return nil, currentTimeRes{}, err
	}

	t, err := util.ToTime(data.Time, "America/New_York")
	if err != nil {
		return nil, currentTimeRes{}, err
	}

	return nil, currentTimeRes{
		CurrentTime: t,
	}, nil
}

type getCurrentConditionsReq struct{}

type getCurrentConditionsRes struct {
	CurrentConditions *weather.CurrentConditions `json:"current_conditions" jsonschema:"Current weather conditions in the local area"`
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
		CurrentConditions: &data.CurrentConditions,
	}, nil
}

type hourlyForecastReq struct{}

type hourlyForecastRes struct {
	Hours []weather.HourlyForecast `json:"hours" jsonschema:"Hourly weather forecast for the local area for the next 48 hours"`
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
		Hours: data.HourlyForecast,
	}, nil
}

type dailyForecastReq struct{}

type dailyForecastRes struct {
	Days []weather.DailyForecast `json:"days" jsonschema:"Daily weather forecast for the local area for the next 10 days"`
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
		Days: data.DailyForecast,
	}, nil
}
