package jq

import (
	"context"
	"encoding/json"
	"fmt"
	"kellegous/jqmcp/internal/weather"
	"time"

	"github.com/itchyny/gojq"
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
	Time string `json:"time" jsonschema:"Current time in local time zone in ISO 8601 format"`
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
	return nil, currentTimeRes{
		Time: data.Time.Format(time.RFC3339),
	}, nil
}

type getCurrentConditionsReq struct {
	JQExpression string `json:"jq_expression,omitempty" jsonschema:"A JQ expression to transform current_conditions property in the result to select only the data that is needed."`
}

type getCurrentConditionsRes struct {
	FilteredCurrentConditions any                        `json:"filtered_current_conditions,omitempty" jsonschema:"The current_conditions property filtered by the JQ expression"`
	CurrentConditions         *weather.CurrentConditions `json:"current_conditions,omitempty" jsonschema:"Current weather conditions in the local area"`
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

	if e := args.JQExpression; e != "" {
		filtered, err := jqRun(ctx, e, data.CurrentConditions)
		if err != nil {
			return nil, getCurrentConditionsRes{}, err
		}
		return nil, getCurrentConditionsRes{
			FilteredCurrentConditions: filtered,
		}, nil
	}

	return nil, getCurrentConditionsRes{
		CurrentConditions: &data.CurrentConditions,
	}, nil
}

type hourlyForecastReq struct {
	JQExpression string `json:"jq_expression,omitempty" jsonschema:"A JQ expression to transform the hours property in the result to select only the data that is needed."`
}

type hourlyForecastRes struct {
	FilteredHours any                      `json:"filtered_hours,omitempty" jsonschema:"The hours property filtered by the JQ expression"`
	Hours         []weather.HourlyForecast `json:"hours,omitempty" jsonschema:"Hourly weather forecast for the local area for the next 48 hours"`
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

	if e := args.JQExpression; e != "" {
		filtered, err := jqRun(ctx, e, data.HourlyForecast)
		if err != nil {
			return nil, hourlyForecastRes{}, err
		}
		return nil, hourlyForecastRes{
			FilteredHours: filtered,
		}, nil
	}

	return nil, hourlyForecastRes{
		Hours: data.HourlyForecast,
	}, nil
}

type dailyForecastReq struct {
	JQExpression string `json:"jq_expression,omitempty" jsonschema:"A JQ expression to transform the days property in the result to select only the data that is needed."`
}

type dailyForecastRes struct {
	FilteredDays any                     `json:"filtered_days,omitempty" jsonschema:"The days property filtered by the JQ expression"`
	Days         []weather.DailyForecast `json:"days,omitempty" jsonschema:"Daily weather forecast for the local area for the next 10 days"`
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

	if e := args.JQExpression; e != "" {
		filtered, err := jqRun(ctx, e, data.DailyForecast)
		if err != nil {
			return nil, dailyForecastRes{}, err
		}
		return nil, dailyForecastRes{
			FilteredDays: filtered,
		}, nil
	}

	return nil, dailyForecastRes{
		Days: data.DailyForecast,
	}, nil
}

func toPrimitive(data any) (any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func jqRun(ctx context.Context, expr string, data any) (any, error) {
	query, err := gojq.Parse(expr)
	if err != nil {
		return nil, err
	}

	raw, err := toPrimitive(data)
	if err != nil {
		return nil, err
	}

	iter := query.RunWithContext(ctx, raw)

	var values []any
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if _, ok := v.(error); ok {
			return nil, fmt.Errorf("jq error: %w", v.(error))
		}

		values = append(values, v)
	}

	if len(values) == 1 {
		return values[0], nil
	}

	return values, nil
}
