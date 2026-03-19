package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/kellegous/poop"
)

var (
	dimColor   = color.New(color.Faint)
	greenColor = color.New(color.FgGreen)
)

type toolCall struct {
	ID   string
	Name string
	Args json.RawMessage
}

type outputStream struct {
	w                io.Writer
	thinking         bool
	pendingToolCalls map[string]*toolCall
}

func newOutputStream(w io.Writer) *outputStream {
	return &outputStream{
		w:                w,
		pendingToolCalls: make(map[string]*toolCall),
	}
}

func (s *outputStream) writeEvent(evt agentEvent) error {
	switch evt := evt.(type) {
	case *ToolCallEvent:
		return s.writeToolCall(evt)
	case *AssistantMessageChunkEvent:
		if evt.IsThinking {
			s.thinking = true
			if _, err := dimColor.Fprint(s.w, evt.Chunk); err != nil {
				return poop.Chain(err)
			}
		} else {
			if s.thinking {
				s.thinking = false
				if _, err := fmt.Fprintln(s.w); err != nil {
					return poop.Chain(err)
				}
			}

			if _, err := fmt.Fprint(s.w, evt.Chunk); err != nil {
				return poop.Chain(err)
			}
		}
		return nil
	case *AssistantDoneEvent:
		if _, err := fmt.Fprintln(s.w); err != nil {
			return poop.Chain(err)
		}
	}
	return nil
}

func (s *outputStream) writeToolCall(evt *ToolCallEvent) error {
	if _, err := greenColor.Fprintf(s.w, "%s(%s)\n", evt.ToolName, evt.ToolArgs.Format()); err != nil {
		return poop.Chain(err)
	}

	if _, err := greenColor.Fprintln(s.w, evt.ToolResult.Format()); err != nil {
		return poop.Chain(err)
	}

	return nil
}
