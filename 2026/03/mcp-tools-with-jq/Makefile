all: bin/agent

bin/agent: cmd/agent/main.go $(shell find internal -type f)
	go build -o $@ $<

clean:
	rm -rf bin

.PHONY: all clean