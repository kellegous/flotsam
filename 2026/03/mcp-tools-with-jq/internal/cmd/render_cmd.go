package cmd

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var mdConverter = goldmark.New(
	goldmark.WithExtensions(extension.Table),
)

type Event struct {
	Message    string          `json:"message,omitempty"`
	Thinking   string          `json:"thinking,omitempty"`
	Content    string          `json:"content,omitempty"`
	ToolName   string          `json:"tool_name,omitempty"`
	ToolArgs   json.RawMessage `json:"tool_args,omitempty"`
	ToolID     string          `json:"tool_id,omitempty"`
	ToolResult json.RawMessage `json:"tool_result,omitempty"`
}

var renderCmd = func() *cobra.Command {
	return &cobra.Command{
		Use:   "render",
		Short: "Pre-render agent content to markdown in agent log files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRender(args[0])
		},
	}
}()

func runRender(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return poop.Chain(err)
	}

	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return poop.Chain(err)
	}

	for i := range events {
		if t := events[i].Thinking; t != "" {
			h, convErr := markdownToHTML(t)
			if convErr != nil {
				return poop.Chain(convErr)
			}
			events[i].Thinking = h
		}

		if c := events[i].Content; c != "" {
			h, convErr := markdownToHTML(c)
			if convErr != nil {
				return poop.Chain(convErr)
			}
			events[i].Content = h
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(events); err != nil {
		return poop.Chain(err)
	}
	return nil
}

func markdownToHTML(src string) (string, error) {
	if src == "" {
		return "", nil
	}
	var buf bytes.Buffer
	if err := mdConverter.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
