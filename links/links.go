package links

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Linker struct {
	Users map[string]*User
	Links map[string]*ShortLink
}

func NewLinker() *Linker {
	l := new(Linker)
	l.ReloadLinks()
	return l
}

type ShortLink struct {
	Short       string     `json:"short,omitempty"`
	Original    string     `json:"original,omitempty"`
	Owner       *User      `json:"owner,omitempty"`
	DateCreated time.Time  `json:"date_created"`
	DateExpired *time.Time `json:"date_expired,omitempty"`
}

type Role = int

const (
	RoleUser Role = iota
	RoleAdmin
)

type User struct {
	Username string
	Email    string
	Role     Role
}

func (l *Linker) GetLink(short string) *ShortLink {
	if l.Links == nil {
		l.ReloadLinks()
	}
	return l.Links[short]
}

func NewLink(original string, opts ...func(*ShortLink)) ShortLink {
	link := ShortLink{
		Original:    original,
		DateCreated: time.Now().UTC(),
	}

	for _, f := range opts {
		f(&link)
	}

	return link
}

func WithShort(short string) func(*ShortLink) {
	if short == "" {
		short = generateShortLink()
	}
	return func(link *ShortLink) {
		link.Short = short
	}
}

func WithRandomShort() func(*ShortLink) {
	return func(link *ShortLink) {
		link.Short = generateShortLink()
	}
}

func WithOwner(owner User) func(*ShortLink) {
	return func(link *ShortLink) {
		link.Owner = &owner
	}
}

func (l *Linker) AddLink(link ShortLink) error {
	if l.Links == nil {
		l.ReloadLinks()
	}
	if link.Short == "" {
		link.Short = generateShortLink()
	}

	l.Links[link.Short] = &link

	return l.SaveLinks()
}

func generateShortLink() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (l *Linker) SaveLinks() error {
	f, err := os.Create("links.json")
	if err != nil {
		return err
	}
	defer f.Close()

	shorts := make([]*ShortLink, 0, len(l.Links))
	for _, link := range l.Links {
		shorts = append(shorts, link)
	}

	return json.NewEncoder(f).Encode(shorts)
}

func (l *Linker) ReloadLinks() {
	var err error
	l.Links, err = loadLinks()
	if err != nil {
		l.Links = make(map[string]*ShortLink)
	}

	l.Users = map[string]*User{
		"admin": {
			Username: "admin",
			Email:    "admin@mail.com",
			Role:     RoleAdmin,
		},
	}
}

func loadLinks() (map[string]*ShortLink, error) {
	_, err := os.Stat("links.json")
	if err != nil {
		return nil, createLinksFile()
	}

	f, err := os.Open("links.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var shorts []*ShortLink
	err = json.NewDecoder(f).Decode(&shorts)
	if err != nil {
		return nil, err
	}

	linksMap := make(map[string]*ShortLink)
	for _, link := range shorts {
		linksMap[link.Short] = link
	}

	return linksMap, nil
}

func createLinksFile() error {
	f, err := os.Create("links.json")
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode([]ShortLink{})
}
