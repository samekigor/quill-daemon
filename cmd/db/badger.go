package db

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/samekigor/quill-daemon/cmd/utils"
)

var dbPath string
var db *badger.DB

func ChangeDbPath(newPath string) (err error) {
	if newPath == "" {
		log.Fatalf("ChangeDbPath: received empty path")
		return err
	}

	dbPath = newPath
	err = utils.CreateDir(newPath)
	if err != nil {
		log.Fatalf("Failed to initialize DB directory: %v", err)
		return err
	}
	return nil
}

func InitDb() (err error) {
	dbPath = "/var/lib/quill/db/"
	if err = ChangeDbPath(dbPath); err != nil {
		utils.ErrorLogger.Fatalf("Db init failure %v", err)
		return err
	}
	opts := badger.DefaultOptions(dbPath).WithLoggingLevel(badger.WARNING)

	db, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	utils.InfoLogger.Printf("Database works in: %s", dbPath)
	return nil
}
