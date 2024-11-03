package db

import (
	"fmt"
	"link-shortener/links"
	"time"
)

const (
	selectUserRole = `
	SELECT role
	FROM users
	WHERE username = ?
    `

	selectUser = `
	SELECT username, email, role
	FROM users
	WHERE username = ?
	`

	selectLink = `
	SELECT short, original, owner, date_created, date_expired
	FROM links
	WHERE short = ?
	`
)

func (db *Sqlite) IsAdmin(username string) (bool, error) {
	role, err := db.SelectUserRole(username)
	if err != nil {
		return false, err
	}
	return role == links.RoleAdmin, nil
}

func (db *Sqlite) SelectUserRole(username string) (links.Role, error) {
	var role links.Role
	err := db.QueryRow(selectUserRole, username).Scan(&role)
	if err != nil {
		return 0, err
	}
	return role, nil
}

func (db *Sqlite) SelectUser(username string) (*links.User, error) {
	var user links.User
	err := db.QueryRow(selectUser, username).Scan(&user.Username, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Sqlite) SelectLink(short string) (*links.ShortLink, error) {
	var link links.ShortLink
	var ownerString, dateCreated *string
	err := db.QueryRow(selectLink, short).Scan(
		&link.Short, &link.Original,
		&ownerString,
		&dateCreated, &link.DateExpired,
	)
	if err != nil {
		return nil, err
	}

	if ownerString != nil {
		role, err := db.SelectUser(*ownerString)
		if err != nil {
			return nil, err
		}
		link.Owner = role
	}

	if dateCreated != nil {
		link.DateCreated, err = time.Parse(time.RFC3339, *dateCreated)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date_created: %w", err)
		}
	}

	return &link, nil
}
