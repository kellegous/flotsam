package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kellegous/poop"
	"trpc.group/trpc-go/trpc-agent-go/event"
)

type toolCall struct {
	ID   string
	Name string
	Args json.RawMessage
}

type outputStream struct {
	w                io.Writer
	thinkingShowing  bool
	pendingToolCalls map[string]*toolCall
}

func newOutputStream(w io.Writer) *outputStream {
	return &outputStream{
		w:                w,
		pendingToolCalls: make(map[string]*toolCall),
	}
}

func (s *outputStream) processEvent(ctx context.Context, evt *event.Event) error {
	switch evt.Object {
	case "chat.completion.chunk":
		return s.processChatCompletionChunk(ctx, evt)
	case "chat.completion":
		return s.processChatCompletion(ctx, evt)
	case "tool.response":
		return s.processToolResponse(ctx, evt)
	}
	return nil
}

func (s *outputStream) reset() {
	s.thinkingShowing = false
}

func (s *outputStream) processChatCompletionChunk(ctx context.Context, evt *event.Event) error {
	res := evt.Response
	if res == nil {
		return poop.New("chat.completion.chunk has no response")
	}

	if len(res.Choices) == 0 {
		return poop.New("chat.completion.chunk has no choices")
	}

	choice := res.Choices[0]

	isThinking := choice.Delta.ReasoningContent != ""
	if isThinking {
		if s.thinkingShowing {
			return nil
		}

		if _, err := fmt.Fprintln(s.w, "Thinking..."); err != nil {
			return poop.Chain(err)
		}
		s.thinkingShowing = true
		return nil
	}

	if _, err := fmt.Fprint(s.w, choice.Delta.Content); err != nil {
		return poop.Chain(err)
	}

	return nil
}

func (s *outputStream) processChatCompletion(ctx context.Context, evt *event.Event) error {
	s.reset()

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

func (s *outputStream) processToolResponse(ctx context.Context, evt *event.Event) error {
	s.reset()

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

		if _, err := fmt.Fprintf(s.w, "%s(%s)\n", call.Name, string(args)); err != nil {
			return poop.Chain(err)
		}

		res, err := json.MarshalIndent(json.RawMessage(msg.Content), "", "  ")
		if err != nil {
			return poop.Chain(err)
		}

		if _, err := fmt.Fprintln(s.w, string(res)); err != nil {
			return poop.Chain(err)
		}

		delete(s.pendingToolCalls, msg.ToolID)
	}

	return nil
}
