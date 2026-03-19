package cmd

import (
	"context"
	"encoding/json"
	"iter"

	"github.com/kellegous/poop"
	"trpc.group/trpc-go/trpc-agent-go/event"
)

type agentEvent interface {
	isAgentEvent()
}

type UserMessageEvent struct {
	Message string
}

func (UserMessageEvent) isAgentEvent() {}

type AssistantMessageChunkEvent struct {
	IsThinking bool
	Chunk      string
}

func (AssistantMessageChunkEvent) isAgentEvent() {}

type AssistantDoneEvent struct{}

func (AssistantDoneEvent) isAgentEvent() {}

type ToolCallEvent struct {
	ToolName   string
	ToolArgs   maybeJSON
	ToolID     string
	ToolResult maybeJSON
}

func (ToolCallEvent) isAgentEvent() {}

func toEvents(ctx context.Context, ch <-chan *event.Event) iter.Seq2[agentEvent, error] {
	return func(yield func(agentEvent, error) bool) {
		pendingToolCalls := make(map[string]*ToolCallEvent)

		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-ch:
				if !ok {
					return
				}

				switch e.Object {
				case "chat.completion.chunk":
					evt, err := processChatCompletionChunk(e)
					if err != nil {
						yield(nil, poop.Chain(err))
						return
					}
					if !yield(evt, nil) {
						return
					}
				case "chat.completion":
					events, err := processChatCompletion(e)
					if err != nil {
						yield(nil, poop.Chain(err))
						return
					}

					for _, evt := range events {
						pendingToolCalls[evt.ToolID] = evt
					}

					if !yield(&AssistantDoneEvent{}, nil) {
						return
					}
				case "tool.response":
					events, err := processToolResponse(e, pendingToolCalls)
					if err != nil {
						yield(nil, poop.Chain(err))
						return
					}
					for _, evt := range events {
						if !yield(evt, nil) {
							return
						}
					}
				}
			}
		}
	}
}

func processChatCompletion(e *event.Event) ([]*ToolCallEvent, error) {
	res := e.Response
	if res == nil {
		return nil, poop.New("chat.completion has no response")
	}

	if len(res.Choices) == 0 {
		return nil, poop.New("chat.completion has no choices")
	}

	choice := res.Choices[0]

	toolCalls := choice.Message.ToolCalls
	if len(toolCalls) == 0 {
		return nil, nil
	}

	events := make([]*ToolCallEvent, 0, len(toolCalls))
	for _, call := range toolCalls {
		events = append(events, &ToolCallEvent{
			ToolName: call.Function.Name,
			ToolID:   call.ID,
			ToolArgs: maybeJSON(call.Function.Arguments),
		})
	}

	return events, nil
}

func processChatCompletionChunk(e *event.Event) (*AssistantMessageChunkEvent, error) {
	res := e.Response
	if res == nil {
		return nil, poop.New("chat.completion.chunk has no response")
	}

	if len(res.Choices) == 0 {
		return nil, poop.New("chat.completion.chunk has no choices")
	}

	choice := res.Choices[0]

	if choice.Delta.ReasoningContent != "" {
		return &AssistantMessageChunkEvent{
			IsThinking: true,
			Chunk:      choice.Delta.ReasoningContent,
		}, nil
	}

	return &AssistantMessageChunkEvent{
		IsThinking: false,
		Chunk:      choice.Delta.Content,
	}, nil
}

func processToolResponse(
	e *event.Event,
	pendingToolCalls map[string]*ToolCallEvent,
) ([]*ToolCallEvent, error) {
	res := e.Response
	if res == nil {
		return nil, poop.New("tool.response has no response")
	}

	events := make([]*ToolCallEvent, 0, len(res.Choices))
	for _, choice := range res.Choices {
		msg := choice.Message
		if msg.Role != "tool" {
			return nil, poop.New("tool.response has no tool message")
		}

		call, ok := pendingToolCalls[msg.ToolID]
		if !ok {
			return nil, poop.New("tool.response has no pending tool call")
		}

		var err error
		call.ToolResult, err = decodeToolResult(msg.Content)
		if err != nil {
			return nil, poop.Chain(err)
		}

		events = append(events, call)
		delete(pendingToolCalls, msg.ToolID)
	}

	return events, nil
}

func decodeToolResult(content string) (maybeJSON, error) {
	var msg []struct {
		Text string `json:"text"`
	}

	if err := json.Unmarshal([]byte(content), &msg); err != nil {
		return nil, poop.Chain(err)
	}

	if len(msg) != 1 {
		return nil, poop.New("tool.response has multiple messages")
	}

	return maybeJSON(msg[0].Text), nil
}

type maybeJSON []byte

func (m maybeJSON) Format() string {
	var raw json.RawMessage
	if err := json.Unmarshal(m, &raw); err != nil {
		return string(m)
	}

	b, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return string(m)
	}

	return string(b)
}

type fallibleJSON struct {
	Valid json.RawMessage
	Raw   string
}

func (f *fallibleJSON) UnmarshalJSON(b []byte) error {
	var valid json.RawMessage
	if err := json.Unmarshal(b, &valid); err != nil {
		f.Raw = string(b)
		return nil
	}
	f.Valid = valid
	return nil
}

func (f fallibleJSON) MarshalJSON() ([]byte, error) {
	if f.Valid != nil {
		return f.Valid, nil
	}
	return json.Marshal(f.Raw)
}

func (f *fallibleJSON) String() string {
	return f.Raw
}
