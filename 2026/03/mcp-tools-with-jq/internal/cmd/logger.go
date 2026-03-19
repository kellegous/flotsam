package cmd

import (
	"encoding/json"
	"io"
)

type logger struct {
	events []agentEvent
}

func (l *logger) writeEvent(evt agentEvent) {
	l.events = append(l.events, evt)
}

func (l *logger) writeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(l.events)
}
