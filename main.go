package main

import (
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tailscale.com/tsnet"
)

var (
	tailscaleHost      = flag.String("tailscale-hostname", "tsp", "hostname for tailscale device")
	tailscalePort      = flag.Int("tailscale-port", 80, "port to listen on for incoming connections from tailscale")
	tailscaleConfigDir = flag.String("tailscale-config-dir", "config", "path to store tailscale configuration")
	tailscaleAuthKey   = flag.String("tailscale-auth-key", "", "tailscale auth key for connecting to the network. If blank, interactive auth will be required")
	upstream           = flag.String("upstream", "", "ip:port of the upstream service to proxy connections to")
	logLevel           = flag.String("log-level", "info", "level of logs to output")
)

func main() {
	envflag.Parse()

	level := slog.Level(0)
	_ = level.UnmarshalText([]byte(*logLevel))

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	_, _, err := net.SplitHostPort(*upstream)
	if err != nil {
		logger.Error("Error parsing upstream address", "error", err)
		os.Exit(1)
	}

	serv := tsnet.Server{
		Hostname: *tailscaleHost,
		Dir:      *tailscaleConfigDir,
		AuthKey:  *tailscaleAuthKey,
		UserLogf: slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo).Printf,
		Logf:     slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug).Printf,
	}

	listener, err := serv.Listen("tcp", fmt.Sprintf(":%d", *tailscalePort))
	if err != nil {
		logger.Error("Error listening on tailnet", "port", *tailscalePort, "error", err)
		os.Exit(1)
	}

	logger.Info("Listening for incoming connections", "hostname", *tailscaleHost, "port", *tailscalePort)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Warn("Error accepting connection", "error", err)
				continue
			}

			go proxy(conn)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("Shutting down...")

	listener.Close()
	serv.Close()
}

func proxy(connection net.Conn) {
	defer connection.Close()
	logger := slog.Default().With("peer", connection.RemoteAddr())

	logger.Debug("Connection accepted")

	up, err := net.Dial("tcp", *upstream)
	if err != nil {
		logger.Debug("Error connecting to upstream", "error", err)
		return
	}
	defer up.Close()

	go func() {
		_, err = io.Copy(up, connection)
		if err != nil {
			logger.Debug("Error copying from upstream to client", "error", err)
		}
	}()

	_, err = io.Copy(connection, up)
	if err != nil {
		logger.Debug("Error copying from client to upstream", "error", err)
	}

	logger.Debug("Connection closed")
}
