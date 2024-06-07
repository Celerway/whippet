package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.golang/paho"
	"github.com/perbu/whippet/whippet"
)

//go:embed .version
var embeddedVersion string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	err := run(ctx, os.Stderr, os.Stdout, os.Stdin, os.Args[1:], os.Environ())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logoutput, output io.Writer, input io.Reader, args, env []string) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	logHandle := slog.NewTextHandler(logoutput, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandle)
	config, showhelp, err := whippet.GetConfig(args)
	if err != nil {
		return fmt.Errorf("getConfig: %w", err)
	}
	if showhelp {
		return nil
	}
	logger.Info("client starting up", "version", embeddedVersion)
	client, msgChan, err := whippet.Connect(runCtx, config, logger)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	go func() {
		<-runCtx.Done()
		logger.Info("context cancelled, disconnecting client")
		if client != nil {
			d := &paho.Disconnect{ReasonCode: 0}
			_ = client.Disconnect(d)
		}
		// Just assume the rest of the application is able to shut down
	}()

	// read the input until OEF
	payload, err := read(ctx, input)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	response, err := whippet.Request(runCtx, client, config, payload, msgChan, logger)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	_, err = output.Write(response)
	if err != nil {
		return fmt.Errorf("output.Write: %w", err)
	}
	return nil
}

func read(ctx context.Context, input io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	done := make(chan error, 1)

	go func() {
		_, err := io.Copy(&buf, input)
		done <- err
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("io.Copy: %w", err)
		}
		return buf.Bytes(), nil
	}
}
