package cmd

import (
	"fmt"
	"iter"
	"strings"

	"github.com/kellegous/poop"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/anthropic"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
)

var defaultModel = func() Model {
	mod, err := parseModel("")
	if err != nil {
		poop.HitFan(err)
	}
	return Model{
		Model: mod,
		val:   "local:",
	}
}()

type Model struct {
	model.Model
	val string
}

func (m *Model) Set(v string) error {
	mod, err := parseModel(v)
	if err != nil {
		return err
	}
	m.val = v
	m.Model = mod
	return nil
}

func (m *Model) String() string {
	return m.val
}

func (m *Model) Type() string {
	return "provider:options"
}

func parseModel(v string) (model.Model, error) {
	if v == "" {
		return fromLocal("")
	}

	prov, opts, ok := strings.Cut(v, ":")
	if !ok {
		return nil, fmt.Errorf("invalid model %q: missing ':'", v)
	}

	opts = strings.TrimSpace(opts)
	switch strings.TrimSpace(prov) {
	case "anthropic":
		return fromAnthropic(opts)
	case "openai":
		return fromOpenAI(opts)
	case "local":
		return fromLocal(opts)
	}
	return nil, fmt.Errorf("invalid model %q: unknown provider %q", v, prov)
}

func fromAnthropic(v string) (model.Model, error) {
	opts := struct {
		model  string
		apiKey string
	}{
		model: "claude-sonnet-4-6",
	}

	for opt, err := range parseOptions(v) {
		if err != nil {
			return nil, poop.Chain(err)
		}
		switch opt.key {
		case "model":
			opts.model = opt.val
		case "api-key":
			opts.apiKey = opt.val
		default:
			return nil, poop.Newf("anthropic:invalid option %q: unknown key %q", opts, opt.key)
		}
	}

	if opts.apiKey == "" {
		return nil, poop.New("anthropic:apiKey is required")
	}

	return anthropic.New(opts.model, anthropic.WithAPIKey(opts.apiKey)), nil
}

func fromOpenAI(v string) (model.Model, error) {
	cfg := struct {
		model   string
		apiKey  string
		baseURL string
	}{
		model: "gpt-4o-mini",
	}

	for opt, err := range parseOptions(v) {
		if err != nil {
			return nil, poop.Chain(err)
		}
		switch opt.key {
		case "model":
			cfg.model = opt.val
		case "api-key":
			cfg.apiKey = opt.val
		case "base-url":
			cfg.baseURL = opt.val
		default:
			return nil, poop.Newf("openai:invalid option %q: unknown key %q", cfg, opt.key)
		}
	}

	if cfg.apiKey == "" {
		return nil, poop.New("openai:apiKey is required")
	}

	opts := []openai.Option{
		openai.WithAPIKey(cfg.apiKey),
	}

	if cfg.baseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.baseURL))
	}

	return openai.New(cfg.model, opts...), nil
}

func fromLocal(v string) (model.Model, error) {
	cfg := struct {
		model   string
		apiKey  string
		baseURL string
	}{
		model:   "gpt-oss:20b",
		apiKey:  "ollama",
		baseURL: "http://127.0.0.1:11434/v1",
	}

	for opt, err := range parseOptions(v) {
		if err != nil {
			return nil, poop.Chain(err)
		}
		switch opt.key {
		case "model":
			cfg.model = opt.val
		case "api-key":
			cfg.apiKey = opt.val
		case "base-url":
			cfg.baseURL = opt.val
		}
	}

	return openai.New(cfg.model, openai.WithAPIKey(cfg.apiKey), openai.WithBaseURL(cfg.baseURL)), nil
}

type option struct {
	key string
	val string
}

func parseOptions(opts string) iter.Seq2[option, error] {
	return func(yield func(option, error) bool) {
		opts := strings.TrimSpace(opts)
		if opts == "" {
			return
		}

		for part := range strings.SplitSeq(opts, ",") {
			part = strings.TrimSpace(part)
			key, val, ok := strings.Cut(part, "=")
			if !ok {
				if !yield(option{}, fmt.Errorf("invalid option %q: missing '='", part)) {
					return
				}
				continue
			}
			if key == "" {
				if !yield(option{}, fmt.Errorf("invalid option %q: key cannot be empty", part)) {
					return
				}
				continue
			}
			if !yield(option{key: key, val: val}, nil) {
				return
			}
		}
	}
}
