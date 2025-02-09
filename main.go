package main

import (
	"github.com/samekigor/quill-daemon/cmd/db"
	"github.com/samekigor/quill-daemon/cmd/server"
	"github.com/samekigor/quill-daemon/cmd/utils"
)

const (
	SocketPath string = "/var/run/quill.sock"
)

func main() {

	utils.InitLogger()
	var err error

	if err = db.InitDb(); err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	if err := server.StartGRPCServer(SocketPath); err != nil {
		utils.ErrorLogger.Fatalf("Failed to start gRPC server: %v", err)
	}

}
