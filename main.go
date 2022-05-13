package main

import (
	"log"

	"github.com/rlawnsxo131/madre-server-v2/database"
	"github.com/rlawnsxo131/madre-server-v2/lib/logger"
	"github.com/rlawnsxo131/madre-server-v2/server"
)

func main() {
	db, err := database.GetDatabaseInstance()
	if err != nil {
		logger.GetDefaultLogger().
			Fatal().Err(err).Msg("")
	}
	defer db.DB.Close()
	// TODO: It should be written to run only in the develop environment.
	database.ExcuteInitSQL(db.DB)

	s := server.New(db)
	s.Start()
}

// TODO: setting up environment
func init() {
	log.Println("init main")
}
