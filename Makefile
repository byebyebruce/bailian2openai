VERSION?=$(shell git describe --tags --always)

build:
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o bailian2openai ./cmd/bailian2openai