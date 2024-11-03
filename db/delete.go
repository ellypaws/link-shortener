package db

import (
	"errors"
	"link-shortener/links"
)

const (
	deleteLink = `
	DELETE FROM links
	WHERE short = ?
	`
)

func (db *Sqlite) DeleteLink(short links.ShortLink) error {
	original, err := db.SelectLink(short.Short)
	if err != nil {
		return err
	}

	if short.Owner.Role != links.RoleAdmin && short.Owner.Username != original.Owner.Username {
		return errors.New("unauthorized")
	}

	_, err = db.Exec(deleteLink, short)
	return err
}
