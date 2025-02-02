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

func SaveToDb(key string, val string) (err error) {
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(key), []byte(val))
		return err
	})
	if err != nil {
		utils.InfoLogger.Printf("Failed to save key %s: %v", key, err)
	}
	return nil
}
