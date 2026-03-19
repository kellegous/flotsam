package cmd

import (
	"encoding/json"
	"io"
)

type logger struct {
	events          []agentEvent
	shouldLogChunks bool
}

func newLogger(shouldLogChunks bool) *logger {
	return &logger{
		events:          make([]agentEvent, 0),
		shouldLogChunks: shouldLogChunks,
	}
}

func (l *logger) writeEvent(evt agentEvent) {
	switch evt := evt.(type) {
	case *AssistantMessageChunkEvent:
		if l.shouldLogChunks {
			l.events = append(l.events, evt)
		}
	default:
		l.events = append(l.events, evt)
	}
}

func (l *logger) writeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(l.events)
}
