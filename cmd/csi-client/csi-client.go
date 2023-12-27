// Package main is the entrypoint for the csi-client.
package main

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"unix:///Users/wave/go/src/csi-driver/csi.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("could not dial unix sock: %s\n", err)
	}

	nodeClient := csi.NewNodeClient(conn)

	resp, err := nodeClient.NodeGetInfo(context.TODO(), &csi.NodeGetInfoRequest{})
	if err != nil {
		fmt.Printf("could not get node info response: %s\n", err)
	}

	fmt.Printf("response: %v\n", resp)
}
