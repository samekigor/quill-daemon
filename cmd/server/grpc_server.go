package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/samekigor/quill-daemon/cmd/utils"
	"github.com/samekigor/quill-daemon/proto/auth"
)

// StartGRPCServer uruchamia serwer gRPC nasłuchujący na Unix Socket.
func StartGRPCServer(socketPath string) error {
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err == nil {
		rootUID := 0
		group, err := user.LookupGroup("quill")
		if err != nil {
			utils.ErrorLogger.Fatalf("failed to get group ID: %v", err)
		}
		quillGID, _ := strconv.Atoi(group.Gid)
		if err := os.Chown(socketPath, rootUID, quillGID); err != nil {
			return fmt.Errorf("failed to set ownership on socket %s: %w", socketPath, err)
		}
		if err := os.Chmod(socketPath, 0660); err != nil {
			return fmt.Errorf("failed to set permissions on socket %s: %v", socketPath, err)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to listen on socket %s: %w", socketPath, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, &AuthServer{})

	log.Printf("gRPC server is listening on Unix Socket: %s", socketPath)
	return grpcServer.Serve(listener)
}

func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
