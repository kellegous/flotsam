package agent

import (
	"context"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/mcp"
)

const (
	App     = "meteorologist"
	User    = "kellegous"
	Session = "9f29f412b8a6a8849527357b813807ce9a1538af"
)

type Runner interface {
	Run(ctx context.Context, message string, opts ...agent.RunOption) (<-chan *event.Event, error)
	Close() error
}

type runnerImpl struct {
	user    string
	session string
	impl    runner.Runner
}

func (r *runnerImpl) Run(
	ctx context.Context,
	message string,
	opts ...agent.RunOption,
) (<-chan *event.Event, error) {
	return r.impl.Run(
		ctx,
		r.user,
		r.session,
		model.NewUserMessage(message),
		opts...,
	)
}

func (r *runnerImpl) Close() error {
	return r.impl.Close()
}

func New(
	ctx context.Context,
	mcpURL string,
	m model.Model,
) Runner {
	mcpTools := mcp.NewMCPToolSet(
		mcp.ConnectionConfig{
			Transport: "streamable",
			ServerURL: mcpURL,
			Timeout:   10 * time.Second,
		},
	)

	sessions := inmemory.NewSessionService()
	agent := llmagent.New(
		"meteorologist",
		llmagent.WithModel(m),
		llmagent.WithGenerationConfig(
			model.GenerationConfig{
				Stream: true,
			},
		),
		llmagent.WithInstruction("Add instruction"),
		llmagent.WithToolSets([]tool.ToolSet{mcpTools}),
	)

	return &runnerImpl{
		user:    User,
		session: Session,
		impl:    runner.NewRunner(App, agent, runner.WithSessionService(sessions)),
	}
}
