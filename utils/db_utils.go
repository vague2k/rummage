package utils

import (
	"database/sql"
	"fmt"
)

// Creates a "items" table if it does not exist in the Rummage DB.
func CreateDBTable(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            entry TEXT,
            score FLOAT,
            lastAccessed INTEGER
        )
    `)
	if err != nil {
		msg := fmt.Sprintf("Could not create 'items' table in rummage db: \n%s", err)
		logger.Fatal(msg)
	}
}
