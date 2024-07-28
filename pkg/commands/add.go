package commands

import "github.com/vague2k/rummage/pkg/database"

func Add(db *database.RummageDB, args ...string) error {
	if _, err := db.AddMultiItems(args...); err != nil {
		return err
	}
	return nil
}
