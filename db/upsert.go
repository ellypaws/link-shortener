package db

import (
	"fmt"
	"link-shortener/links"
	"time"
)

const (
	upsertUser = `
	INSERT INTO users (username, email, password, role)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(username) DO UPDATE SET
		email = excluded.email,
		password = excluded.password,
		role = excluded.role
	`

	upsertLink = `
	INSERT INTO links (short, original, owner, date_created, date_expired)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(short) DO UPDATE SET
		original = excluded.original,
		owner = excluded.owner,
		date_created = excluded.date_created,
		date_expired = excluded.date_expired
	`
)

func (db *Sqlite) UpsertUser(user *links.User) error {
	_, err := db.Exec(
		upsertUser,
		user.Username,
		user.Email,
		user.Role,
	)
	return err
}

func (db *Sqlite) UpsertLink(link *links.ShortLink) error {
	if link == nil {
		return fmt.Errorf("link is nil")
	}

	if link.Owner != nil {
		err := db.UpsertUser(link.Owner)
		if err != nil {
			return fmt.Errorf("failed to upsert user for link: %w", err)
		}

		_, err = db.Exec(
			upsertLink,
			link.Short,
			link.Original,
			link.Owner.Username,
			link.DateCreated.UTC().Format(time.RFC3339),
			nil,
		)
		return fmt.Errorf("failed to upsert link: %w", err)
	}

	_, err := db.Exec(
		upsertLink,
		link.Short,
		link.Original,
		nil,
		link.DateCreated.UTC().Format(time.RFC3339),
		nil,
	)
	return err
}
