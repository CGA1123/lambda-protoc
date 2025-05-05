package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {
	Program   string
	Arguments []string
}

func (h *Handler) Invoke(ctx context.Context, payload string) (string, error) {
	stdin := base64.NewDecoder(base64.StdEncoding, strings.NewReader(payload))
	out := &bytes.Buffer{}
	stdout := base64.NewEncoder(base64.StdEncoding, out)

	cmd := exec.Command(h.Program, h.Arguments...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running %s %s: %w", h.Program, strings.Join(h.Arguments, " "), err)
	}

	if err := stdout.Close(); err != nil {
		return "", fmt.Errorf("closing stdout base64 encoder: %w", err)
	}

	return out.String(), nil
}

func main() {
	prog := os.Args[1:]
	if len(prog) == 0 {
		os.Exit(1)
	}

	h := Handler{Program: prog[0], Arguments: prog[1:]}

	lambda.Start(h.Invoke)
}
