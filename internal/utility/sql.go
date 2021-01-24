package utility

import (
	"database/sql"
	"log"
	"time"
)

func DBAlive(db *sql.DB) {
	log.Println("connecting to database... ")
	for {
		// Ping by itself is un-reliable, the connections are cached. This
		// ensures that the database is still running by executing a harmless
		// dummy query against it.
		_, err := db.Exec("SELECT true")
		if err == nil {
			return
		}
		log.Println("retrying...")
		time.Sleep(time.Second)
	}
}

