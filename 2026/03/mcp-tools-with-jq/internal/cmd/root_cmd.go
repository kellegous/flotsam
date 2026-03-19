package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"kellegous/jqmcp/internal/agent"
	"kellegous/jqmcp/internal/mcp/jq"
	"kellegous/jqmcp/internal/mcp/plain"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/kellegous/poop"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	mcpAddr string
	model   Model
	useJQ   bool
	logFile string
}

var rootCmd = func() *cobra.Command {
	flags := rootFlags{
		model: defaultModel,
	}

	rootCmd := &cobra.Command{
		Use:   "jqmcp",
		Short: "jqmcp is a test of MCP tools accepting JQ expressions",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runRoot(cmd.Context(), &flags); err != nil {
				poop.HitFan(err)
			}
		},
	}

	rootCmd.Flags().StringVar(
		&flags.mcpAddr,
		"http.addr",
		"localhost:4000",
		"Address of the http server that serves MCP endpoints",
	)

	rootCmd.Flags().Var(
		&flags.model,
		"model",
		"Model to use for the agent",
	)

	rootCmd.Flags().BoolVar(
		&flags.useJQ,
		"use-jq",
		false,
		"Use JQ to filter the MCP responses",
	)

	rootCmd.Flags().StringVar(
		&flags.logFile,
		"log-file",
		"",
		"File to write the log to",
	)

	return rootCmd
}()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func urlForAddr(addr string, useJQ bool) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if host == "" {
		host = "localhost"
	}

	if useJQ {
		return fmt.Sprintf("http://%s:%s/jq", host, port), nil
	}

	return fmt.Sprintf("http://%s:%s/plain", host, port), nil
}

func runRoot(ctx context.Context, flags *rootFlags) error {
	mcpURL, err := urlForAddr(flags.mcpAddr, flags.useJQ)
	if err != nil {
		return poop.Chain(err)
	}

	var lg logger

	ch := make(chan error)

	go func() {
		ch <- serveHTTP(ctx, flags.mcpAddr)
	}()

	go func() {
		ch <- runAgent(
			ctx,
			agent.New(ctx, mcpURL, flags.model.Model),
			&lg,
		)
	}()

	if err := <-ch; err != nil {
		return poop.Chain(err)
	}

	if flags.logFile != "" {
		w, err := os.Create(flags.logFile)
		if err != nil {
			return poop.Chain(err)
		}
		defer w.Close()

		if err := lg.writeTo(w); err != nil {
			return poop.Chain(err)
		}
	}

	return nil
}

type discard struct {
	io.Writer
}

func (d *discard) Close() error {
	return nil
}

func openLogFile(path string) (io.WriteCloser, error) {
	if path == "" {
		return &discard{Writer: io.Discard}, nil
	}
	return os.Create(path)
}

func runAgent(
	ctx context.Context,
	r agent.Runner,
	lg *logger,
) error {
	scanner := bufio.NewScanner(os.Stdin)

	stream := newOutputStream(os.Stdout)

	for {
		fmt.Fprint(os.Stdout, "> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "":
			continue
		case "/exit", "/quit":
			return nil
		}

		events, err := r.Run(ctx, input)
		if err != nil {
			return poop.Chain(err)
		}

		for evt, err := range toEvents(ctx, events) {
			if err != nil {
				return poop.Chain(err)
			}

			if err := stream.writeEvent(evt); err != nil {
				return poop.Chain(err)
			}

			lg.writeEvent(evt)
		}
	}
	return scanner.Err()
}

func serveHTTP(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	plainSrv := plain.New(ctx)
	mux.Handle("/plain", http.StripPrefix("/plain", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return plainSrv
	}, nil)))

	jqSrv := jq.New(ctx)
	mux.Handle("/jq", http.StripPrefix("/jq", mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return jqSrv
	}, nil)))

	return http.ListenAndServe(addr, mux)
}
