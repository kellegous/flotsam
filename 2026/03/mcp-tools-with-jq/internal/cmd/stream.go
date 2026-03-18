package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/kellegous/poop"
	"trpc.group/trpc-go/trpc-agent-go/event"
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
	hadThinking      bool
	pendingToolCalls map[string]*toolCall
}

func newOutputStream(w io.Writer) *outputStream {
	return &outputStream{
		w:                w,
		pendingToolCalls: make(map[string]*toolCall),
	}
}

func (s *outputStream) processEvent(evt *event.Event) error {
	switch evt.Object {
	case "chat.completion.chunk":
		return s.processChatCompletionChunk(evt)
	case "chat.completion":
		return s.processChatCompletion(evt)
	case "tool.response":
		return s.processToolResponse(evt)
	}
	return nil
}

func (s *outputStream) processChatCompletionChunk(evt *event.Event) error {
	res := evt.Response
	if res == nil {
		return poop.New("chat.completion.chunk has no response")
	}

	if len(res.Choices) == 0 {
		return poop.New("chat.completion.chunk has no choices")
	}

	choice := res.Choices[0]

	if choice.Delta.ReasoningContent != "" {
		s.hadThinking = true
		if _, err := dimColor.Fprint(s.w, choice.Delta.ReasoningContent); err != nil {
			return poop.Chain(err)
		}
		return nil
	}

	if s.hadThinking {
		s.hadThinking = false
		if _, err := fmt.Fprintln(s.w); err != nil {
			return poop.Chain(err)
		}
	}

	if _, err := fmt.Fprint(s.w, choice.Delta.Content); err != nil {
		return poop.Chain(err)
	}

	return nil
}

func (s *outputStream) processChatCompletion(evt *event.Event) error {
	if _, err := fmt.Fprintln(s.w); err != nil {
		return poop.Chain(err)
	}

	res := evt.Response
	if res == nil {
		return poop.New("chat.completion has no response")
	}

	if len(res.Choices) == 0 {
		return poop.New("chat.completion has no choices")
	}

	choice := res.Choices[0]

	toolCalls := choice.Message.ToolCalls
	if len(toolCalls) == 0 {
		return nil
	}

	for _, call := range toolCalls {
		s.pendingToolCalls[call.ID] = &toolCall{
			ID:   call.ID,
			Name: call.Function.Name,
			Args: call.Function.Arguments,
		}
	}

	return nil
}

func (s *outputStream) processToolResponse(evt *event.Event) error {
	res := evt.Response
	if res == nil {
		return poop.New("tool.response has no response")
	}

	for _, choice := range res.Choices {
		msg := choice.Message
		if msg.Role != "tool" {
			return poop.New("tool.response has no tool message")
		}

		call, ok := s.pendingToolCalls[msg.ToolID]
		if !ok {
			return poop.New("tool.response has no pending tool call")
		}

		args, err := json.MarshalIndent(call.Args, "", "  ")
		if err != nil {
			return poop.Chain(err)
		}

		if _, err := greenColor.Fprintf(s.w, "%s(%s)\n", call.Name, string(args)); err != nil {
			return poop.Chain(err)
		}

		res, err := json.MarshalIndent(json.RawMessage(msg.Content), "", "  ")
		if err != nil {
			return poop.Chain(err)
		}

		if _, err := greenColor.Fprintln(s.w, string(res)); err != nil {
			return poop.Chain(err)
		}

		delete(s.pendingToolCalls, msg.ToolID)
	}

	return nil
}
