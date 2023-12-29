// Package main is the entrypoint for the csi-client.
package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type envConfig struct {
	SocketPath string `env:"SOCKET_PATH"`
}

func main() {
	logger := zap.Must(zap.NewDevelopment())

	envVars := &envConfig{}

	err := env.Parse(envVars)
	if err != nil {
		log.Fatalf("failed to connect to parse env vars: %s", err)
	}

	conn, err := grpc.Dial(
		envVars.SocketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to socket: %s", err)
	}

	nodeClient := csi.NewNodeClient(conn)

	resp, err := nodeClient.NodeGetInfo(context.TODO(), &csi.NodeGetInfoRequest{})
	if err != nil {
		log.Fatalf("could not get node info response: %s", err)
	}

	logger.Info("successful request",
		zap.Any("response", resp),
	)
}
